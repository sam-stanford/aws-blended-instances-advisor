package api

import (
	"ec2-test/aws/types"
	"ec2-test/cache"
	"ec2-test/config"
	"ec2-test/utils"
	"testing"
	"time"
)

const (
	CONFIG_FILEPATH          = "../../../../config.json"
	CONFIG_API_DOWNLOADS_DIR = "../../../../assets/downloads/"
	CACHE_FILEPATH           = "../../../../assets/test/test-api-cache"
	MAX_INSTANCES            = 100
	REGION1                  = "us-west-1"
	REGION2                  = "ca-central-1" // TODO: Why no on-demand for "eu-west-2"
)

func getTestRegions() (*types.Region, *types.Region, error) {
	r1, err := types.NewRegion(REGION1)
	if err != nil {
		return nil, nil, err
	}
	r2, err := types.NewRegion(REGION2)
	if err != nil {
		return nil, nil, err
	}
	return &r1, &r2, nil
}

func TestGetInstances(t *testing.T) {
	cfg, err := config.ParseConfig(CONFIG_FILEPATH)
	if err != nil {
		t.Fatalf("Failed to read config: %s", err.Error())
	}
	cfg.ApiConfig.MaxInstancesToFetch = MAX_INSTANCES
	cfg.ApiConfig.DownloadsDir = CONFIG_API_DOWNLOADS_DIR

	logger, err := utils.CreateMockLogger()
	if err != nil {
		t.Fatalf("Failed to create mock logger: %s", err.Error())
	}

	cache, err := cache.New(CACHE_FILEPATH)
	if err != nil {
		t.Fatalf("Failed to create cache: %s", err.Error())
	}

	region1, region2, err := getTestRegions()
	if err != nil {
		t.Fatalf("Error occured when parsing test regions: %s", err.Error())
	}

	noCacheStartTime := time.Now()
	regionInstanceMap, err := GetInstances(&cfg.ApiConfig, []types.Region{*region1, *region2}, &cfg.Credentials.Test, cache, logger)
	if err != nil {
		t.Fatalf("Error thrown when fecthing instances: %s", err.Error())
	}
	noCacheEndTime := time.Now()

	instancesRegion1, ok := regionInstanceMap[*region1]
	if !ok || len(instancesRegion1) == 0 {
		t.Fatalf("Zero instances returned for region %s", REGION1)
	}
	if len(instancesRegion1) > cfg.ApiConfig.MaxInstancesToFetch*2 {
		t.Fatalf(
			"More instances returned than config max for region %s. Wanted: < %d, got: %d",
			REGION1,
			config.DEFAULT_API_MAX_INSTANCES_TO_FETCH,
			len(instancesRegion1),
		)
	}

	instancesRegion2, ok := regionInstanceMap[*region2]
	if !ok || len(instancesRegion2) == 0 {
		t.Fatalf("Zero instances returned for region %s", REGION2)
	}
	if len(instancesRegion2) > cfg.ApiConfig.MaxInstancesToFetch*2 {
		t.Fatalf(
			"More instances returned than config max for region %s. Wanted: < %d, got: %d",
			REGION2,
			config.DEFAULT_API_MAX_INSTANCES_TO_FETCH,
			len(instancesRegion2),
		)
	}

	cachedStartTime := time.Now()
	cachedRegionInstanceMap, err := GetInstances(&cfg.ApiConfig, []types.Region{*region1, *region2}, &cfg.Credentials.Test, cache, logger)
	if err != nil {
		t.Fatalf("Error thrown when fecthing instances for a second time: %s", err.Error())
	}
	cachedEndTime := time.Now()

	cachedInstancesRegion1, ok := cachedRegionInstanceMap[*region1]
	if !ok || len(cachedInstancesRegion1) == 0 {
		t.Fatalf("Zero instances returned for region %s", REGION1)
	}
	if len(cachedInstancesRegion1) > cfg.ApiConfig.MaxInstancesToFetch*2 {
		t.Fatalf(
			"More instances returned than config max for region %s. Wanted: < %d, got: %d",
			REGION1,
			config.DEFAULT_API_MAX_INSTANCES_TO_FETCH,
			len(cachedInstancesRegion1),
		)
	}
	if len(cachedInstancesRegion1) != len(instancesRegion1) {
		t.Fatalf(
			"Different number of returned from cache than request. Wanted: %d, got: %d",
			len(instancesRegion1),
			len(cachedInstancesRegion1),
		)
	}

	cachedInstancesRegion2, ok := cachedRegionInstanceMap[*region2]
	if !ok || len(cachedInstancesRegion2) == 0 {
		t.Fatalf("Zero instances returned for region %s", REGION2)
	}
	if len(cachedInstancesRegion2) > cfg.ApiConfig.MaxInstancesToFetch*2 {
		t.Fatalf(
			"More instances returned than config max for region %s. Wanted: < %d, got: %d",
			REGION2,
			config.DEFAULT_API_MAX_INSTANCES_TO_FETCH,
			len(cachedInstancesRegion2),
		)
	}
	if len(cachedInstancesRegion2) != len(instancesRegion2) {
		t.Fatalf(
			"Different number of returned from cache than request. Wanted: %d, got: %d",
			len(instancesRegion2),
			len(cachedInstancesRegion2),
		)
	}

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
