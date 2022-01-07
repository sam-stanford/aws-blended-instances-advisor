package api

import (
	types "aws-blended-instances-advisor/aws/types"
	"aws-blended-instances-advisor/config"
	instPkg "aws-blended-instances-advisor/instances"
	"context"

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

// GetOnDemandInstances fetches on-demand instance offerings from the
// AWS API, returning them as a list of Instances.
//
// An error is returned if a critical failure is encountered during
// the processes execution, with handleable failures being logged and
// handled appropriately.
func GetOnDemandInstances(
	config *config.AwsApiConfig,
	regions []types.Region,
	creds credentials.StaticCredentialsProvider,
	logger *zap.Logger,
) (
	map[types.Region][]*instPkg.Instance,
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
	cfg *config.AwsApiConfig,
	pricingClient *pricing.Client,
	regions []types.Region,
	maxInstanceCount int,
	logger *zap.Logger,
) (
	map[types.Region][]*instPkg.Instance,
	error,
) {

	regionToInstancesMap := make(map[types.Region][]*instPkg.Instance)

	for _, region := range regions {
		regionInstances := make([]*instPkg.Instance, 0)

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
				zap.String("region", region.CodeString()),
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
			"finished fetching on-demand instances for region",
			zap.String("region", region.CodeString()),
			zap.Int("totalInstanceCount", total),
			zap.Int("maxInstanceCount", total),
		)

		if len(regionInstances) > maxInstanceCount && maxInstanceCount > 0 {
			logger.Info(
				"removed excess instances to keep to max instance count",
				zap.String("region", region.CodeString()),
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
	locationFilterValue := region.NameString()

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
