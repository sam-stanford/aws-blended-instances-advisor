package aws

import (
	"context"
	"ec2-test/config"
	"ec2-test/utils"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/pricing/types"
	"go.uber.org/zap"
)

// TODO: Move creds into another file

// TODO: Store in config/ fetch these & store / fetch AZs too
type Region string

const (
	UsEast1      Region = "us-east-1"
	UsEast2      Region = "us-east-2"
	UsWest1      Region = "us-west-1"
	UsWest2      Region = "us-west-2"
	AfSouth1     Region = "af-south-1"
	ApEast1      Region = "ap-east-1"
	ApSouth1     Region = "ap-south-1"
	ApNorthEast3 Region = "ap-northeast-3"
	ApNorthEast2 Region = "ap-northeast-2"
	ApSouthEast1 Region = "ap-southeast-1"
	ApSouthEast2 Region = "ap-southeast-2"
	ApNorthEast1 Region = "ap-northeast-1"
	CaCentral1   Region = "ca-central-1"
	EuCentral1   Region = "eu-central-1"
	EuWest1      Region = "eu-west-1"
	EuWest2      Region = "eu-west-2"
	EuSouth1     Region = "eu-south-1"
	EuWest3      Region = "eu-west-3"
	EuNorth1     Region = "eu-north-1"
	MeSouth1     Region = "me-south-1"
	SaEast1      Region = "sa-east-1"
	UsGovEast1   Region = "us-gov-east-1"
	UsGovWest1   Region = "us-gov-west-1"
)

const (
	AWS_PRICING_API_REGION = "us-east-1" // Only us-east-1 works currently
)

type OperatingSystem string

const (
	Linux   OperatingSystem = "linux"
	Windows OperatingSystem = "windows"
	MacOs   OperatingSystem = "macos"
)

type Instance struct {
	Name                  string
	Memory                float64
	Vcpu                  int
	Region                Region
	OperatingSystem       OperatingSystem
	PricePerHour          float64
	RevocationProbability float64
}

func (product *product) productToInstance() (*Instance, error) {
	return &Instance{}, nil // TODO
}

func FetchInstanceInfo(config *config.Config, logger *zap.Logger) ([]Instance, error) {

	creds := createAwsCredentials(config.Credentials.AwsKeyId, config.Credentials.AwsSecretKey)

	spotInstances, err := getSpotInstances(config, creds, logger)
	if err != nil {
		// TODO
		return nil, err
	}
	fmt.Println(spotInstances) // TODO: Remove

	onDemandInstances, err := getOnDemandInstances(config, creds, logger)
	if err != nil {
		// TODO
		return nil, err
	}
	fmt.Println(onDemandInstances) // TODO: Remove

	return nil, nil
}

func createAwsCredentials(keyId string, secretKey string) credentials.StaticCredentialsProvider {
	return credentials.NewStaticCredentialsProvider(keyId, secretKey, "")
}

func getSpotInstances(config *config.Config, creds credentials.StaticCredentialsProvider, logger *zap.Logger) ([]Instance, error) {
	// for _, region := range config.Constraints.Regions {
	// 	getSpotInstancesForRegion(region, config, creds, logger)
	// }

	fetchSpotInstanceRevocationInfo(config, logger)

	return nil, nil // TODO
}

func getSpotInstancesForRegion(region string, config *config.Config, creds credentials.StaticCredentialsProvider, logger *zap.Logger) ([]Instance, error) {
	awsConfig, err := createAwsConfig(region, creds)
	if err != nil {
		return nil, err
	}

	ec2Client := createEc2Client(awsConfig)
	fetchSpotInstanceAvailabilityInfo(ec2Client, logger)

	return nil, nil // TODO
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

func fetchSpotInstanceAvailabilityInfo(ec2Client *ec2.Client, logger *zap.Logger) {
	response, err := ec2Client.DescribeSpotPriceHistory(context.TODO(), &ec2.DescribeSpotPriceHistoryInput{})
	if err != nil {
		log.Printf("error: %v", err) // TODO
		return
	}

	fmt.Println(response)

	for _, instance := range response.SpotPriceHistory {
		fmt.Printf(
			"AZ: %s, Type: %v, Desc: %v, Price: %s, Time: %v\n",
			*instance.AvailabilityZone,
			instance.InstanceType,
			instance.ProductDescription,
			*instance.SpotPrice,
			instance.Timestamp,
		)
	}

}

func fetchSpotInstanceRevocationInfo(config *config.Config, logger *zap.Logger) error {
	cwd, err := utils.GetCallerPath()
	if err != nil {
		return err
	}

	filepath, err := utils.CreateFilepath(cwd, config.DownloadsDir, "spot-instance-info.json")
	if err != nil {
		return err
	}

	err = utils.DownloadFile(config.Endpoints.AwsSpotInstanceInfoUrl, filepath)
	logger.Info("downloaded spot instance revocation data",
		zap.String("getUrl", config.Endpoints.AwsSpotInstanceInfoUrl),
		zap.String("downloadFilepath", filepath), // TODO: Pass logger to download instead (as need to log cache check)
	)

	return err
}

func createAwsPricingClient(awsCredentials credentials.StaticCredentialsProvider) *pricing.Client {
	return pricing.New(pricing.Options{
		Region:      AWS_PRICING_API_REGION, // TODO: Is this the fixed one? If so, comment required
		Credentials: awsCredentials,
	})
}

func getOnDemandInstances(config *config.Config, creds credentials.StaticCredentialsProvider, logger *zap.Logger) ([]Instance, error) {
	pricingClient := createAwsPricingClient(creds)
	fetchOnDemandInstanceInfo(pricingClient, logger)
	return nil, nil // TODO
}

func fetchOnDemandInstanceInfo(pricingClient *pricing.Client, logger *zap.Logger) {
	// TODO: Return error

	ec2ServiceCode := "AmazonEC2"
	locationFilterName := "location"
	locationFilterValue := "Canada (Central)" // TODO: Use a loop for this val from config

	total := 0

	nextToken := ""
	firstIter := true
	for nextToken != "" || firstIter {

		resp, err := pricingClient.GetProducts(context.TODO(), &pricing.GetProductsInput{
			ServiceCode: &ec2ServiceCode,
			NextToken:   &nextToken,
			Filters:     []types.Filter{{Field: &locationFilterName, Value: &locationFilterValue, Type: "TERM_MATCH"}},
		})
		if err != nil {
			log.Printf("error: %v", err) // TODO: use zap
			return
		}

		total += len(resp.PriceList)

		for _, productString := range resp.PriceList {
			var product product
			err = json.Unmarshal([]byte(productString), &product)
			if err != nil {
				// logger("error: %v", err)
				// TODO: Error
				return
			}
		}

		// fmt.Printf("%+v\n", resp)
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
