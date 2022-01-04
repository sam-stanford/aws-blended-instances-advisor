package schema

import (
	types "aws-blended-instances-advisor/aws/types"
	"aws-blended-instances-advisor/config"
	instPkg "aws-blended-instances-advisor/instances"
	"aws-blended-instances-advisor/utils"
	"context"
	"encoding/json"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"go.uber.org/zap"
)

// TODO: Add log of how many not parsed (like on-demand.go line 66)

func GetSpotInstances(
	config *config.AwsApiConfig,
	regions []types.Region,
	creds credentials.StaticCredentialsProvider,
	logger *zap.Logger,
) (map[types.Region][]*instPkg.Instance, error) {

	regionToInstanceMap := make(map[types.Region][]*instPkg.Instance)

	regionRevocationInfoMap, instanceSpecMap, err := fetchSpotInstanceRevocationInfoAndSpecsMap(config, logger)
	if err != nil {
		logger.Error("error fetching spot instance revocation info and specifications from API", zap.Error(err))
		return nil, err
	}

	for _, region := range regions {
		logger.Info("creating spot instances for region", zap.String("region", region.CodeString()))

		regionRevocationInfo, ok := regionRevocationInfoMap[region.CodeString()]
		if !ok {
			// TODO: Handle more gracefully
			logger.Error("could not find region revocation info", zap.String("region", region.CodeString()))
			continue
		}
		logger.Debug("fetched revocation info")

		regionSpotPrices, err := getSpotInstancePricesForRegion(config, region, creds, logger)
		if err != nil {
			// TODO: Handle more gracefully
			logger.Error(
				"could not fetch region spot prices",
				zap.String("region", region.CodeString()),
				zap.Error(err),
			)
			return nil, err
		}

		regionPriceMap := createInstancePriceMap(regionSpotPrices)

		instances, err := createRegionSpotInstances(config, region, &regionRevocationInfo, regionPriceMap, instanceSpecMap, logger)
		if err != nil {
			return nil, err
		}
		regionToInstanceMap[region] = instances
		logger.Debug("Finished creating instances")
		logger.Info("TODO") // TODO: Log & count
	}

	return regionToInstanceMap, nil
}

func createInstancePriceMap(spotPrices []ec2Types.SpotPrice) map[string]ec2Types.SpotPrice {
	instancePriceMap := make(map[string]ec2Types.SpotPrice)
	for _, price := range spotPrices {
		instancePriceMap[string(price.InstanceType)] = price
	}
	return instancePriceMap
}

func getSpotInstancePricesForRegion(
	config *config.AwsApiConfig,
	region types.Region,
	creds credentials.StaticCredentialsProvider,
	logger *zap.Logger,
) ([]ec2Types.SpotPrice, error) {
	awsConfig, err := createAwsConfig(region.CodeString(), creds)
	if err != nil {
		return nil, err
	}

	ec2Client := createEc2Client(awsConfig)
	logger.Info("created EC2 client")

	return fetchSpotInstanceAvailabilityInfo(ec2Client, config.MaxInstancesToFetch, logger)
}

func fetchSpotInstanceAvailabilityInfo(
	ec2Client *ec2.Client,
	maxInstanceCount int,
	logger *zap.Logger,
) (
	[]ec2Types.SpotPrice,
	error,
) {

	spotPrices := make([]ec2Types.SpotPrice, 0)

	nextToken := ""
	firstIter := true
	total := 0
	for (total < maxInstanceCount || maxInstanceCount <= 0) && (nextToken != "" || firstIter) {
		resp, err := ec2Client.DescribeSpotPriceHistory(context.TODO(), &ec2.DescribeSpotPriceHistoryInput{})
		if err != nil {
			logger.Error("error calling DescribeSpotInstancePriceHistory to EC2 client", zap.Error(err))
			return nil, err
		}
		logger.Info("fetched spot instance prices", zap.Int("count", len(resp.SpotPriceHistory)))

		spotPrices = append(spotPrices, resp.SpotPriceHistory...)
		total += len(spotPrices)

		firstIter = false
		if resp.NextToken != nil {
			nextToken = *resp.NextToken
		} else {
			nextToken = ""
		}
	}

	logger.Info(
		"finished fetching spot instances pricing info",
		zap.Int("totalInstanceCount", len(spotPrices)),
		zap.Int("maxInstanceCount", maxInstanceCount),
	)

	// TODO: Ensure next block's working properly
	if len(spotPrices) > maxInstanceCount {
		logger.Info(
			"removed excess instances to keep to max instance count",
			zap.Int("removed", len(spotPrices)-maxInstanceCount),
		)
		spotPrices = spotPrices[:maxInstanceCount]
	}

	return spotPrices, nil
}

func fetchSpotInstanceRevocationInfoAndSpecsMap(
	config *config.AwsApiConfig,
	logger *zap.Logger,
) (
	map[string]regionSpotInstanceRevocationInfo,
	map[string]spotInstanceSpecs,
	error,
) {
	cwd, err := utils.GetCallerPath()
	if err != nil {
		logger.Error("failed to fetch current working directory", zap.Error(err))
		return nil, nil, err
	}

	filepath, err := utils.CreateFilepath(cwd, config.DownloadsDir, "spot-instance-info.json")
	if err != nil {
		logger.Error("failed to create filepath", zap.Error(err))
		return nil, nil, err
	}

	err = utils.DownloadFile(config.Endpoints.AwsSpotInstanceInfoUrl, filepath) // TODO: Return boolean on whether cache was used
	if err != nil {
		logger.Error("failed to download file", zap.String("getUrl", config.Endpoints.AwsSpotInstanceInfoUrl), zap.Error(err))
		return nil, nil, err
	}
	logger.Info("downloaded spot instance revocation data",
		zap.String("getUrl", config.Endpoints.AwsSpotInstanceInfoUrl),
		zap.String("downloadFilepath", filepath),
	)

	// TODO - Separate this out so we can test this with test file
	infoFile, err := utils.FileToBytes(filepath)
	if err != nil {
		logger.Error("failed to parse file to bytes", zap.String("file", filepath), zap.Error(err))
		return nil, nil, err
	}

	var info spotInstancesInfo
	json.Unmarshal(infoFile, &info)
	return info.RegionPrices, info.SpecsMap, nil
}

func createRegionSpotInstances(
	cfg *config.AwsApiConfig,
	region types.Region,
	regionRevocationInfo *regionSpotInstanceRevocationInfo,
	regionInstancePriceMap map[string]ec2Types.SpotPrice,
	instanceSpecMap map[string]spotInstanceSpecs,
	logger *zap.Logger,
) (
	[]*instPkg.Instance,
	error,
) {
	instances := make([]*instPkg.Instance, 0)

	for instanceType, revocationInfo := range regionRevocationInfo.LinuxInstances {
		spec, ok := instanceSpecMap[instanceType]
		if !ok {
			logger.Debug(
				"failed to create spot instance because no instance specification exists",
				zap.String("instance", instanceType),
			)
			continue
		}

		price, ok := regionInstancePriceMap[instanceType]
		if !ok {
			logger.Debug(
				"failed to create spot instance because no price exists for instance",
				zap.String("instance", instanceType),
			)
			continue
		}

		instance, err := createInstanceFromSpotInstanceInfo(&price, &revocationInfo, &spec, region, "Linux")
		if err != nil {
			logger.Debug("failed to create instance from given spot instance info", zap.Error(err))
			continue
		}
		instances = append(instances, instance)
	}

	for instanceType, revocationInfo := range regionRevocationInfo.WindowsInstances {
		spec, ok := instanceSpecMap[instanceType]
		if !ok {
			logger.Debug(
				"failed to create spot instance because no instance specification exists",
				zap.String("instance", instanceType),
			)
			continue
		}

		price, ok := regionInstancePriceMap[instanceType]
		if !ok {
			logger.Debug("failed to create spot instance")
			continue
		}

		instance, err := createInstanceFromSpotInstanceInfo(&price, &revocationInfo, &spec, region, "Windows")
		if err != nil {
			logger.Debug("failed to create instance from given spot instance info", zap.Error(err))
			continue
		}

		instances = append(instances, instance)
	}

	return instances, nil
}

func createInstanceFromSpotInstanceInfo(
	spotPrice *ec2Types.SpotPrice,
	revocationInfo *spotInstanceRevocationInfo,
	specs *spotInstanceSpecs,
	region types.Region,
	os string,
) (
	*instPkg.Instance,
	error,
) {

	price, err := parseSpotInstancePrice(spotPrice.SpotPrice)
	if err != nil {
		return nil, err
	}
	revocationProbability, err := revocationInfo.getRevocationProbability()
	if err != nil {
		return nil, err
	}

	return &instPkg.Instance{
		Id:                    utils.GenerateUuid(),
		Name:                  string(spotPrice.InstanceType),
		MemoryGb:              specs.MemoryGb,
		Vcpu:                  specs.Vcpu,
		Region:                region,
		OperatingSystem:       os,
		AvailabilityZone:      *spotPrice.AvailabilityZone,
		PricePerHour:          price,
		RevocationProbability: revocationProbability,
	}, nil
}

func parseSpotInstancePrice(price *string) (float64, error) {
	return strconv.ParseFloat(*price, 64)
}
