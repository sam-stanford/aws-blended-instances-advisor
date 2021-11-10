package aws

import (
	"context"
	"ec2-test/config"
	"ec2-test/utils"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	pricingTypes "github.com/aws/aws-sdk-go-v2/service/pricing/types"
	"go.uber.org/zap"
)

const (
	AWS_PRICING_API_REGION = "us-east-1" // Only us-east-1 works currently
)

// TODO: Doc comment
func FetchInstancesInfo(config *config.Config, logger *zap.Logger) (map[Region][]Instance, error) {

	creds := createAwsCredentials(config.Credentials.AwsKeyId, config.Credentials.AwsSecretKey)

	// spotInstances, err := getSpotInstances(config, creds, logger)
	// if err != nil {
	// 	// TODO
	// 	return nil, err
	// }

	onDemandInstances, err := getOnDemandInstances(config, creds, logger)
	if err != nil {
		// TODO
		return nil, err
	}

	// return append(spotInstances, onDemandInstances...), nil
	return onDemandInstances, err
}

func createAwsCredentials(keyId string, secretKey string) credentials.StaticCredentialsProvider {
	return credentials.NewStaticCredentialsProvider(keyId, secretKey, "")
}

func getSpotInstances(config *config.Config, creds credentials.StaticCredentialsProvider, logger *zap.Logger) (map[Region][]Instance, error) {
	for _, region := range config.Constraints.Regions {
		getSpotInstanceAvailabilityInfoForRegion(region, config, creds, logger)
	}
	fetchSpotInstanceRevocationInfo(config, logger)
	return nil, nil // TODO
}

func getSpotInstanceAvailabilityInfoForRegion(region string, config *config.Config, creds credentials.StaticCredentialsProvider, logger *zap.Logger) ([]ec2Types.SpotPrice, error) {
	awsConfig, err := createAwsConfig(region, creds)
	if err != nil {
		return nil, err
	}

	ec2Client := createEc2Client(awsConfig)
	return fetchSpotInstanceAvailabilityInfo(ec2Client, logger)
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

func fetchSpotInstanceRevocationInfo(config *config.Config, logger *zap.Logger) (*spotInstancesInfo, error) {
	cwd, err := utils.GetCallerPath()
	if err != nil {
		return nil, err
	}

	filepath, err := utils.CreateFilepath(cwd, config.DownloadsDir, "spot-instance-info.json")
	if err != nil {
		return nil, err
	}

	err = utils.DownloadFile(config.Endpoints.AwsSpotInstanceInfoUrl, filepath)
	if err != nil {
		return nil, err
	}
	logger.Info("downloaded spot instance revocation data",
		zap.String("getUrl", config.Endpoints.AwsSpotInstanceInfoUrl),
		zap.String("downloadFilepath", filepath), // TODO: Pass logger to download instead (as need to log cache check)
	)

	// TODO - Separate this out so we can test this with test file
	infoFile, err := utils.FileToBytes(filepath)
	if err != nil {
		return nil, err
	}

	var info spotInstancesInfo
	json.Unmarshal(infoFile, &info)
	return &info, nil
}

func createAwsPricingClient(awsCredentials credentials.StaticCredentialsProvider) *pricing.Client {
	return pricing.New(pricing.Options{
		Region:      AWS_PRICING_API_REGION, // TODO: Is this the fixed one? If so, comment required
		Credentials: awsCredentials,
	})
}

func getOnDemandInstances(config *config.Config, creds credentials.StaticCredentialsProvider, logger *zap.Logger) (map[Region][]Instance, error) {
	pricingClient := createAwsPricingClient(creds)
	fetchOnDemandInstanceInfo(pricingClient, logger)
	return nil, nil // TODO
}

func fetchOnDemandInstanceInfo(pricingClient *pricing.Client, logger *zap.Logger) ([]onDemandInstanceInfo, error) {
	// TODO: See if we can download & cache file
	// TODO: Return error
	instances := make([]onDemandInstanceInfo, 0)

	ec2ServiceCode := "AmazonEC2" // TODO: Make const
	locationFilterName := "location"
	locationFilterValue := "Canada (Central)" // TODO: Use a loop for this val from config

	total := 0

	nextToken := ""
	firstIter := true
	for nextToken != "" || firstIter {

		resp, err := pricingClient.GetProducts(context.TODO(), &pricing.GetProductsInput{
			ServiceCode: &ec2ServiceCode,
			NextToken:   &nextToken,
			Filters:     []pricingTypes.Filter{{Field: &locationFilterName, Value: &locationFilterValue, Type: "TERM_MATCH"}},
		})
		if err != nil {
			return nil, err
		}

		total += len(resp.PriceList)

		for _, productString := range resp.PriceList {
			var info onDemandInstanceInfo
			err = json.Unmarshal([]byte(productString), &info)
			if err != nil {
				return nil, err
			}
			if info.Specs.Attributes.MarketOption == "OnDemand" {
				instances = append(instances, info)
			}
		}

		logger.Info("fetched instances", zap.Int("instanceCount", len(resp.PriceList)))

		firstIter = false
		if resp.NextToken != nil {
			nextToken = *resp.NextToken
		} else {
			logger.Info("finished fetching on-demand instances", zap.Int("totalInstanceCount", total))
			nextToken = ""
		}
	}
}
