package api

import (
	"context"
	. "ec2-test/aws/types"
	"ec2-test/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	pricingTypes "github.com/aws/aws-sdk-go-v2/service/pricing/types"
	"go.uber.org/zap"
)

const (
	AWS_PRICING_API_REGION = "us-east-1" // Only us-east-1 works currently (2021-11-11)
)

// TODO: Doc comment
func GetInstances(config *config.Config, logger *zap.Logger) (map[Region][]Instance, error) {

	creds := createAwsCredentials(config.Credentials.AwsKeyId, config.Credentials.AwsSecretKey)

	onDemandInstances, err := GetOnDemandInstances(config, creds, logger)
	if err != nil {
		return nil, err
	}
	spotInstances, err := GetSpotInstances(config, creds, logger)
	if err != nil {
		return nil, err
	}

	instances := joinSpotAndOnDemandInstances(onDemandInstances, spotInstances, config.GetRegions())
	return instances, nil
}

func createAwsCredentials(keyId string, secretKey string) credentials.StaticCredentialsProvider {
	return credentials.NewStaticCredentialsProvider(keyId, secretKey, "")
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
		Region:      AWS_PRICING_API_REGION, // TODO: Is this the fixed one? If so, comment required
		Credentials: awsCredentials,
	})
}

func joinSpotAndOnDemandInstances(
	onDemandInstances map[Region][]Instance,
	spotInstances map[Region][]Instance,
	regions []Region,
) map[Region][]Instance {
	for _, region := range regions {
		onDemandInstances[region] = append(onDemandInstances[region], spotInstances[region]...)
	}
	return onDemandInstances
}

func getOnDemandInstancesFromApi(
	pricingClient *pricing.Client,
	region Region,
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
