package api

import (
	"context"
	types "ec2-test/aws/types"
	"ec2-test/cache"
	"ec2-test/config"
	"ec2-test/instances"
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	pricingTypes "github.com/aws/aws-sdk-go-v2/service/pricing/types"
	"go.uber.org/zap"
)

const (
	AWS_PRICING_API_REGION   = "us-east-1"      // Only us-east-1 works currently (2021-11-11)
	INSTANCES_CACHE_FILENAME = "instances.json" // TODO: Inject
	INSTANCES_CACHE_DURATION = 672              // TODO: Inject
)

// TODO: Doc comment & use go routines to parallelise fetches
func GetInstances(
	apiConfig *config.ApiConfig,
	regions []types.Region,
	creds *config.Credentials,
	cache *cache.Cache,
	logger *zap.Logger,
) (
	map[types.Region][]instances.Instance,
	error,
) {

	instances, err := getInstancesFromCache(INSTANCES_CACHE_FILENAME, cache)
	if err != nil {
		logger.Info("no instances found in cache", zap.String("reason", err.Error()))
	} else {
		logger.Info("instances retrieved from cache")
		return instances, nil
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

	instances = joinSpotAndOnDemandInstances(onDemandInstances, spotInstances, regions)
	err = storeInstancesInCache(instances, INSTANCES_CACHE_FILENAME, cache)
	if err != nil {
		logger.Error("failed to store instances in cache", zap.Error(err))
		return nil, err
	}

	return instances, nil
}

func getInstancesFromCache(instancesCacheFilename string, c *cache.Cache) (map[types.Region][]instances.Instance, error) {
	isValid := c.IsValid(instancesCacheFilename)
	if isValid {
		instancesFileContent, err := c.Get(instancesCacheFilename)
		if err != nil {
			return nil, err
		}
		var instanceToRegionMap map[types.Region][]instances.Instance
		err = json.Unmarshal([]byte(instancesFileContent), &instanceToRegionMap)
		if err != nil {
			return nil, err
		}
		return instanceToRegionMap, nil
	}
	return nil, errors.New("instances not in cache")
}

func storeInstancesInCache(instanceToRegionMap map[types.Region][]instances.Instance, instancesCacheFilename string, c *cache.Cache) error {
	instancesFileContent, err := json.Marshal(instanceToRegionMap)
	if err != nil {
		return err
	}
	err = c.Set(instancesCacheFilename, string(instancesFileContent), time.Hour*INSTANCES_CACHE_DURATION)
	if err != nil {
		return err
	}
	return nil
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

func joinSpotAndOnDemandInstances(
	onDemandInstances map[types.Region][]instances.Instance,
	spotInstances map[types.Region][]instances.Instance,
	regions []types.Region,
) map[types.Region][]instances.Instance {
	for _, region := range regions {
		onDemandInstances[region] = append(onDemandInstances[region], spotInstances[region]...)
	}
	return onDemandInstances
}

func getOnDemandInstancesFromApi(
	pricingClient *pricing.Client,
	region types.Region,
	nextToken string,
) (*pricing.GetProductsOutput, error) {

	serviceCode := EC2_SERVICE_CODE
	locationFilterKey := LOCATION_FILTER_KEY
	locationFilterValue := region.ToNameString()

	return pricingClient.GetProducts(context.TODO(), &pricing.GetProductsInput{
		ServiceCode: &serviceCode,
		NextToken:   &nextToken,
		Filters: []pricingTypes.Filter{{
			Field: &locationFilterKey,
			Value: &locationFilterValue,
			Type:  TERM_MATCH_FILTER_TYPE,
		}},
	})
}
