package api

import (
	types "ec2-test/aws/types"
	"ec2-test/config"
	"ec2-test/instance"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
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
	map[types.Region][]instance.Instance,
	error,
) {
	pricingClient := createAwsPricingClient(creds)
	return getRegionToOnDemandInstancesMap(
		pricingClient,
		regions,
		config.MaxInstancesToFetch,
		logger,
	)
}

func getRegionToOnDemandInstancesMap(
	pricingClient *pricing.Client,
	regions []types.Region,
	maxInstanceCount int,
	logger *zap.Logger,
) (
	map[types.Region][]instance.Instance,
	error,
) {

	regionToInstancesMap := make(map[types.Region][]instance.Instance)

	for _, region := range regions {
		regionInstances := make([]instance.Instance, 0)

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

			parsedInstances := parseOnDemandApiResponseToInstances(resp, logger)

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
