package api

import (
	"context"
	types "ec2-test/aws/types"
	"ec2-test/config"
	instPkg "ec2-test/instances"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	pricingTypes "github.com/aws/aws-sdk-go-v2/service/pricing/types"
	"go.uber.org/zap"
)

const (
	EC2_SERVICE_CODE       = "AmazonEC2"
	LOCATION_FILTER_KEY    = "location"
	TERM_MATCH_FILTER_TYPE = "TERM_MATCH"
)

func GetOnDemandInstances(
	config *config.ApiConfig,
	regions []types.Region,
	creds credentials.StaticCredentialsProvider,
	logger *zap.Logger,
) (
	map[types.Region][]instPkg.Instance,
	error,
) {
	pricingClient := createAwsPricingClient(creds)
	return getRegionToOnDemandInstancesMap(
		config,
		pricingClient,
		regions,
		config.MaxInstancesToFetch,
		logger,
	)
}

func getRegionToOnDemandInstancesMap(
	cfg *config.ApiConfig,
	pricingClient *pricing.Client,
	regions []types.Region,
	maxInstanceCount int,
	logger *zap.Logger,
) (
	map[types.Region][]instPkg.Instance,
	error,
) {

	regionToInstancesMap := make(map[types.Region][]instPkg.Instance)

	for _, region := range regions {
		regionInstances := make([]instPkg.Instance, 0)

		nextToken := ""
		firstIter := true
		total := 0
		for (total < maxInstanceCount || maxInstanceCount <= 0) && (nextToken != "" || firstIter) {

			resp, err := getOnDemandInstancesFromApi(pricingClient, region, nextToken)
			if err != nil {
				logger.Error("error fetching on-demand instances from API", zap.Error(err))
				return nil, err
			}
			logger.Info("fetched on-demand instances", zap.Int("instanceCount", len(resp.PriceList)))

			parsedInstances := parseOnDemandApiResponseToInstances(cfg, resp, logger)

			logger.Info(
				"parsed on-demand instances",
				zap.Int("parsedCount", len(parsedInstances)),
				zap.Int("skippedCount", len(resp.PriceList)-len(parsedInstances)),
			)

			total += len(parsedInstances)

			regionInstances = append(regionInstances, parsedInstances...)

			firstIter = false
			if resp.NextToken != nil {
				nextToken = *resp.NextToken
			} else {
				nextToken = ""
			}
		}

		logger.Info(
			"fetched on-demand instances for region",
			zap.String("region", region.ToCodeString()),
			zap.Int("totalInstanceCount", total),
			zap.Int("maxInstanceCount", total),
		)

		if len(regionInstances) > maxInstanceCount {
			logger.Info(
				"removed excess instances to keep to max instance count",
				zap.String("region", region.ToCodeString()),
				zap.Int("removed", len(regionInstances)-maxInstanceCount),
			)
			regionInstances = regionInstances[:maxInstanceCount]
		}

		regionToInstancesMap[region] = regionInstances
	}

	return regionToInstancesMap, nil
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
