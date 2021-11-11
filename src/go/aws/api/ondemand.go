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

func GetOnDemandInstances(config *config.Config, creds credentials.StaticCredentialsProvider, logger *zap.Logger) (map[Region][]Instance, error) {
	pricingClient := createAwsPricingClient(creds)
	regions := config.GetRegions()
	return getRegionToOnDemandInstancesMap(pricingClient, regions, logger)
}

func getRegionToOnDemandInstancesMap(pricingClient *pricing.Client, regions []Region, logger *zap.Logger) (map[Region][]Instance, error) {
	// TODO: Cache results

	regionToInstancesMap := make(map[Region][]Instance)

	for _, region := range regions {
		regionInstances := make([]Instance, 0)

		nextToken := ""
		firstIter := true
		totalFetched := 0
		for nextToken != "" || firstIter {

			resp, err := getOnDemandInstancesFromApi(pricingClient, region, nextToken)
			if err != nil {
				return nil, err
			}

			parsedInstances, err := parseOnDemandApiResponseToInstances(resp)
			if err != nil {
				return nil, err
			}

			logger.Info("fetched instances", zap.Int("instanceCount", len(resp.PriceList)))
			totalFetched += len(resp.PriceList)

			regionInstances = append(regionInstances, parsedInstances...)

			firstIter = false
			if resp.NextToken != nil {
				nextToken = *resp.NextToken
			} else {
				nextToken = ""
			}
		}

		regionToInstancesMap[region] = regionInstances
		logger.Info("Fetched on-demand instances for region", zap.String("region", region.ToCodeString()), zap.Int("count", totalFetched))
	}

	return regionToInstancesMap, nil
}
