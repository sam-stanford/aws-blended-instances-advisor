package api

import (
	"context"
	types "ec2-test/aws/types"
	"ec2-test/cache"
	"ec2-test/config"
	instPkg "ec2-test/instances"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"go.uber.org/zap"
)

const (
	AWS_PRICING_API_REGION   = "us-east-1"      // Only us-east-1 works currently (2021-11-11)
	INSTANCES_CACHE_FILENAME = "instances.json" // TODO: Inject
	INSTANCES_CACHE_DURATION = 672              // TODO: Inject
)

// TODO: Doc comment & use go routines to parallelise fetches

func GetInstancesRegionInfoMap(
	apiConfig *config.AwsApiConfig,
	regions []types.Region,
	creds *config.Credentials,
	cache *cache.Cache,
	logger *zap.Logger,
) (
	instPkg.RegionInfoMap,
	error,
) {

	regionInfoMap, err := getRegionInfoMapFromCache(INSTANCES_CACHE_FILENAME, cache)
	if err != nil {
		logger.Info("no instances found in cache", zap.String("reason", err.Error()))
	} else {
		// TODO: Func wrapper for count log & validate fetched instances (or maybe do this in a test )
		for _, region := range regions {
			logger.Info(
				"instances fetched from cache",
				zap.String("region", region.ToCodeString()),
				zap.Int("allInstancesCount", len(regionInfoMap[region].AllInstances.Instances)),
				zap.Int("permanentInstanceCount", len(regionInfoMap[region].PermanentInstances.Instances)),
			)
		}
		return regionInfoMap, nil
	}

	awsCreds := createAwsCredentials(creds)

	onDemandInstances, err := GetOnDemandInstances(apiConfig, regions, awsCreds, logger)
	if err != nil {
		logger.Error("error fetching on-demand instances", zap.Error(err))
		return nil, err
	}
	spotInstances, err := GetSpotInstances(apiConfig, regions, awsCreds, logger)
	if err != nil {
		logger.Error("error fetching spot instances", zap.Error(err))
		return nil, err
	}

	regionInfoMap = createRegionInfoMap(onDemandInstances, spotInstances, regions)
	err = storeRegionInfoMapInCache(regionInfoMap, INSTANCES_CACHE_FILENAME, cache)
	if err != nil {
		logger.Error("failed to store instances in cache", zap.Error(err))
		return nil, err
	}

	// TODO: Func wrapper for count log
	for _, region := range regions {
		logger.Info(
			"stored instances in cache",
			zap.String("region", region.ToCodeString()),
			zap.Int("allInstancesCount", len(regionInfoMap[region].AllInstances.Instances)),
			zap.Int("permanentInstanceCount", len(regionInfoMap[region].PermanentInstances.Instances)),
		)
	}

	return regionInfoMap, nil
}

func createAwsCredentials(creds *config.Credentials) credentials.StaticCredentialsProvider {
	return credentials.NewStaticCredentialsProvider(creds.AwsKeyId, creds.AwsSecretKey, "")
}

func createAwsConfig(awsRegion string, creds credentials.StaticCredentialsProvider) (aws.Config, error) {
	return awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithCredentialsProvider(creds),
		awsConfig.WithRegion(awsRegion),
	)
}

func createEc2Client(awsConfig aws.Config) *ec2.Client {
	return ec2.NewFromConfig(awsConfig)
}

func createAwsPricingClient(awsCredentials credentials.StaticCredentialsProvider) *pricing.Client {
	return pricing.New(pricing.Options{
		Region:      AWS_PRICING_API_REGION,
		Credentials: awsCredentials,
	})
}

func createRegionInfoMap(
	onDemandInstances map[types.Region][]instPkg.Instance,
	spotInstances map[types.Region][]instPkg.Instance,
	regions []types.Region,
) instPkg.RegionInfoMap {

	regionInfoMap := make(instPkg.RegionInfoMap)
	for _, region := range regions {

		thisRegionInfo := instPkg.CreateRegionInfo(onDemandInstances[region], spotInstances[region])
		regionInfoMap[region] = thisRegionInfo
	}
	return regionInfoMap
}
