//go:build integration
// +build integration

package api

import (
	"aws-blended-instances-advisor/aws/types"
	"aws-blended-instances-advisor/cache"
	"aws-blended-instances-advisor/config"
	instPkg "aws-blended-instances-advisor/instances"
	"aws-blended-instances-advisor/utils"
	"errors"
	"fmt"
	"testing"
	"time"
)

const (
	CONFIG_FILEPATH          = "../../../../config.json"
	CONFIG_API_DOWNLOADS_DIR = "testdata/downloads"
	CACHE_FILEPATH           = "testdata/cache"
	MAX_INSTANCES            = 100
)

func TestGetInstances(t *testing.T) {
	cfg, err := config.ParseConfig(CONFIG_FILEPATH)
	if err != nil {
		t.Fatalf("Failed to read config: %s", err.Error())
	}
	cfg.AwsApiConfig.MaxInstancesToFetch = MAX_INSTANCES
	cfg.AwsApiConfig.DownloadsDir = CONFIG_API_DOWNLOADS_DIR

	logger, err := utils.CreateMockLogger()
	if err != nil {
		t.Fatalf("Failed to create mock logger: %s", err.Error())
	}

	cache, err := cache.New(CACHE_FILEPATH)
	if err != nil {
		t.Fatalf("Failed to create cache: %s", err.Error())
	}

	noCacheStartTime := time.Now()
	globalInfo, err := GetInstancesAndInfo(
		&cfg.AwsApiConfig,
		&cfg.Credentials,
		cache,
		logger,
	)
	if err != nil {
		t.Fatalf("Error thrown when fecthing instances: %s", err.Error())
	}
	noCacheEndTime := time.Now()

	err = validateGlobalInfo(globalInfo, 2*cfg.AwsApiConfig.MaxInstancesToFetch)
	if err != nil {
		t.Fatal(err.Error())
	}

	cachedStartTime := time.Now()
	cachedGlobalInfo, err := GetInstancesAndInfo(
		&cfg.AwsApiConfig,
		&cfg.Credentials,
		cache,
		logger,
	)
	if err != nil {
		t.Fatalf("Error thrown when fecthing instances for a second time: %s", err.Error())
	}
	cachedEndTime := time.Now()

	validateCachedGlobalInfo(
		globalInfo,
		cachedGlobalInfo,
		2*cfg.AwsApiConfig.MaxInstancesToFetch,
	)

	noCacheFetchTime := noCacheEndTime.Sub(noCacheStartTime)
	cacheFetchTime := cachedEndTime.Sub(cachedStartTime)

	differenceThreshold := time.Second * 2
	maxNoCacheFetchTime := noCacheFetchTime + differenceThreshold

	if maxNoCacheFetchTime < noCacheFetchTime {
		t.Fatalf(
			"Cached fetch was not significantly faster than no-cache fetch. Wanted: < %v - %v, got: %v",
			noCacheFetchTime,
			differenceThreshold,
			cacheFetchTime,
		)
	}
}

func validateGlobalInfo(
	globalInfo *instPkg.GlobalInfo,
	maxTotalInstances int,
) error {

	if globalInfo.GlobalAggregates.Count == 0 {
		return errors.New("global aggregates not calculated correctly")
	}

	for _, region := range []types.Region{types.UsEast1, types.EuWest2} {

		info, ok := globalInfo.RegionInfoMap[region]

		if !ok || len(info.PermanentInstances) == 0 || len(info.TransientInstances) == 0 {
			return fmt.Errorf("zero instances returned for region %s", region.CodeString())
		}

		if len(info.PermanentInstances) > maxTotalInstances {
			return fmt.Errorf(
				"more permanent instances fetched than max wanted. Fetched: %d, max: %d",
				len(info.PermanentInstances),
				maxTotalInstances,
			)
		}

		if len(info.TransientInstances) > maxTotalInstances {
			return fmt.Errorf(
				"more transient instances fetched than max wanted. Fetched: %d, max: %d",
				len(info.TransientInstances),
				maxTotalInstances,
			)
		}
	}

	return nil
}

func validateCachedGlobalInfo(
	fetchedGlobalInfo *instPkg.GlobalInfo,
	cachedGlobalInfo *instPkg.GlobalInfo,
	maxTotalInstances int,
) error {

	validateGlobalInfo(fetchedGlobalInfo, maxTotalInstances)

	for _, region := range []types.Region{types.UsEast1, types.EuWest2} {
		storedInfo := cachedGlobalInfo.RegionInfoMap[region]
		fetchedInfo := fetchedGlobalInfo.RegionInfoMap[region]

		if len(storedInfo.PermanentInstances) != len(fetchedInfo.PermanentInstances) {
			return fmt.Errorf(
				"Different number of permanent instances fetched from cache than stored. "+
					"Region: %s, stored: %d, fetched: %d",
				region.CodeString(),
				len(storedInfo.PermanentInstances),
				len(fetchedInfo.PermanentInstances),
			)
		}

		if len(storedInfo.TransientInstances) != len(fetchedInfo.TransientInstances) {
			return fmt.Errorf(
				"Different number of transient instances fetched from cache than stored. "+
					"Region: %s, stored: %d, fetched: %d",
				region.CodeString(),
				len(storedInfo.TransientInstances),
				len(fetchedInfo.TransientInstances),
			)
		}

	}

	return nil
}
