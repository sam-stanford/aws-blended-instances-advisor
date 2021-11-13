package api

import (
	. "ec2-test/aws/types"
	"ec2-test/config"

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
	config *config.Config,
	creds credentials.StaticCredentialsProvider,
	logger *zap.Logger,
) (
	map[Region][]Instance,
	error,
) {
	pricingClient := createAwsPricingClient(creds)
	regions := config.GetRegions()
	return getRegionToOnDemandInstancesMap(pricingClient, regions, config.Constraints.MaxInstanceCount, logger)
}

func getRegionToOnDemandInstancesMap(
	pricingClient *pricing.Client,
	regions []Region,
	maxInstanceCount int,
	logger *zap.Logger,
) (
	map[Region][]Instance,
	error,
) {
	// TODO: Cache results - maybe cache entire list of instances instead

	regionToInstancesMap := make(map[Region][]Instance)

	for _, region := range regions {
		regionInstances := make([]Instance, 0)

		nextToken := ""
		firstIter := true
		total := 0
		for (total <= maxInstanceCount || maxInstanceCount <= 0) && (nextToken != "" || firstIter) {

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

		regionToInstancesMap[region] = regionInstances
		logger.Info(
			"fetched on-demand instances for region",
			zap.String("region", region.ToCodeString()),
			zap.Int("totalInstanceCount", total),
			zap.Int("maxInstanceCount", total),
		)
	}

	return regionToInstancesMap, nil
}
