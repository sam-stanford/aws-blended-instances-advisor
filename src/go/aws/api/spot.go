package api

import (
	"context"
	. "ec2-test/aws/types"
	"ec2-test/config"
	"ec2-test/utils"
	"encoding/json"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"go.uber.org/zap"
)

func GetSpotInstances(
	config *config.Config,
	creds credentials.StaticCredentialsProvider,
	logger *zap.Logger,
) (map[Region][]Instance, error) {

	regionToInstanceMap := make(map[Region][]Instance)

	regionRevocationInfoMap, instanceSpecMap, err := fetchSpotInstanceRevocationInfoAndSpecsMap(config, logger)
	if err != nil {
		return nil, err
	}

	for _, region := range config.GetRegions() {

		regionRevocationInfo := regionRevocationInfoMap[region.ToCodeString()]
		if err != nil {
			// TODO: Handle more gracefully
			return nil, err
		}
		regionSpotPrices, err := getSpotInstancePricesForRegion(region, config, creds, logger)
		if err != nil {
			// TODO: Handle more gracefully
			return nil, err
		}
		regionPriceMap := createInstancePriceMap(regionSpotPrices)

		instances, err := createRegionSpotInstances(region, &regionRevocationInfo, regionPriceMap, instanceSpecMap, logger)
		if err != nil {
			return nil, err
		}
		regionToInstanceMap[region] = instances
		logger.Info("TOOD") // TODO: Log & count
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
	region Region,
	config *config.Config,
	creds credentials.StaticCredentialsProvider,
	logger *zap.Logger,
) ([]ec2Types.SpotPrice, error) {
	awsConfig, err := createAwsConfig(region.ToCodeString(), creds)
	if err != nil {
		return nil, err
	}
	ec2Client := createEc2Client(awsConfig)
	return fetchSpotInstanceAvailabilityInfo(ec2Client, logger)
}

func fetchSpotInstanceAvailabilityInfo(ec2Client *ec2.Client, logger *zap.Logger) ([]ec2Types.SpotPrice, error) {
	spotPrices := make([]ec2Types.SpotPrice, 0)

	nextToken := ""
	firstIter := true
	for nextToken != "" || firstIter {
		resp, err := ec2Client.DescribeSpotPriceHistory(context.TODO(), &ec2.DescribeSpotPriceHistoryInput{})
		if err != nil {
			return nil, err
		}

		spotPrices = append(spotPrices, resp.SpotPriceHistory...)

		firstIter = false
		if resp.NextToken != nil {
			nextToken = *resp.NextToken
		} else {
			nextToken = ""
		}
	}

	logger.Info("finished fetching spot instances pricing info", zap.Int("totalInstanceCount", len(spotPrices)))
	return spotPrices, nil
}

func fetchSpotInstanceRevocationInfoAndSpecsMap(
	config *config.Config,
	logger *zap.Logger,
) (
	map[string]regionSpotInstanceRevocationInfo,
	map[string]spotInstanceSpecs,
	error,
) {
	cwd, err := utils.GetCallerPath()
	if err != nil {
		return nil, nil, err
	}

	filepath, err := utils.CreateFilepath(cwd, config.DownloadsDir, "spot-instance-info.json")
	if err != nil {
		return nil, nil, err
	}

	err = utils.DownloadFile(config.Endpoints.AwsSpotInstanceInfoUrl, filepath)
	if err != nil {
		return nil, nil, err
	}
	logger.Info("downloaded spot instance revocation data",
		zap.String("getUrl", config.Endpoints.AwsSpotInstanceInfoUrl),
		zap.String("downloadFilepath", filepath), // TODO: Pass logger to download instead (as need to log cache check)
	)

	// TODO - Separate this out so we can test this with test file
	infoFile, err := utils.FileToBytes(filepath)
	if err != nil {
		return nil, nil, err
	}

	var info spotInstancesInfo
	json.Unmarshal(infoFile, &info)
	return info.RegionPrices, info.SpecsMap, nil
}

func createRegionSpotInstances(
	region Region,
	regionRevocationInfo *regionSpotInstanceRevocationInfo,
	regionInstancePriceMap map[string]ec2Types.SpotPrice,
	instanceSpecMap map[string]spotInstanceSpecs,
	logger *zap.Logger,
) (
	[]Instance,
	error,
) {
	instances := make([]Instance, 0)

	// TODO: Wrap this in method to avoid repeated code for Windows
	for instanceType, revocationInfo := range regionRevocationInfo.LinuxInstances {
		spec, okSpec := instanceSpecMap[instanceType]
		price, okPrice := regionInstancePriceMap[instanceType]
		if okSpec && okPrice {
			instance, err := createInstanceFromSpotInstanceInfo(&price, &revocationInfo, &spec, region, Linux)
			if err != nil {
				// TODO: Do something
			}
			instances = append(instances, *instance)
		} else {
			logger.Info("cannot create spot instance") // TODO: increase verbosity
		}
	}

	for instanceType, revocationInfo := range regionRevocationInfo.WindowsInstances {
		spec, okSpec := instanceSpecMap[instanceType]
		price, okPrice := regionInstancePriceMap[instanceType]
		if okSpec && okPrice {
			instance, err := createInstanceFromSpotInstanceInfo(&price, &revocationInfo, &spec, region, Windows)
			if err != nil {
				// TODO: Do something
			}
			instances = append(instances, *instance)
		} else {
			logger.Info("cannot create spot instance") // TODO: increase verbosity
		}
	}

	return instances, nil
}

func createInstanceFromSpotInstanceInfo(
	spotPrice *ec2Types.SpotPrice,
	revocationInfo *spotInstanceRevocationInfo,
	specs *spotInstanceSpecs,
	region Region,
	os OperatingSystem,
) (
	*Instance,
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

	return &Instance{
		Name:                  string(spotPrice.InstanceType),
		MemoryGb:              specs.MemoryGb,
		Vcpus:                 specs.Vcpus,
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
