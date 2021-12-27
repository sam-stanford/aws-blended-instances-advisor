package schema

import (
	types "aws-blended-instances-advisor/aws/types"
	"aws-blended-instances-advisor/cache"
	"aws-blended-instances-advisor/config"
	instPkg "aws-blended-instances-advisor/instances"
	"context"

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

func GetInstancesAndInfo(
	apiConfig *config.AwsApiConfig,
	creds *config.Credentials,
	cache *cache.Cache,
	logger *zap.Logger,
) (
	*instPkg.GlobalInfo,
	error,
) {

	regions := types.GetAllRegions()

	globalInstanceInfo, err := getGlobalInstanceInfoFromCache(INSTANCES_CACHE_FILENAME, cache)
	if err != nil {
		logger.Info("no instances found in cache", zap.String("reason", err.Error()))
	} else {
		err = globalInstanceInfo.Validate()
		if err == nil {
			globalInstanceInfo.Log("instances and info fetched from cache", logger)
			return globalInstanceInfo, nil
		}
		logger.Warn("invalid instances cache", zap.Error(err))
	}

	logger.Info("fetching instances from AWS API")

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

	globalInfo := instPkg.CreateGlobalInfo(onDemandInstances, spotInstances, regions)

	err = storeGlobalInstanceInfoInCache(globalInfo, INSTANCES_CACHE_FILENAME, cache)
	if err != nil {
		logger.Error("failed to store instances in cache", zap.Error(err))
		return nil, err
	}

	globalInfo.Log("stored instances in cache", logger)

	return &globalInfo, nil
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
