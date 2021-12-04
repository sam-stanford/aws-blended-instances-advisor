//go:build integration
// +build integration

package api

import (
	"ec2-test/aws/types"
	"ec2-test/cache"
	"ec2-test/config"
	instPkg "ec2-test/instances"
	"ec2-test/utils"
	"fmt"
	"testing"
	"time"
)

// TODO: Move to testdata directory
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
	regions := []types.Region{*region1, *region2}

	noCacheStartTime := time.Now()
	regionInfoMap, err := GetInstancesRegionInfoMap(
		&cfg.ApiConfig,
		regions,
		&cfg.Credentials.Test,
		cache,
		logger,
	)
	if err != nil {
		t.Fatalf("Error thrown when fecthing instances: %s", err.Error())
	}
	noCacheEndTime := time.Now()

	err = validateRegionInfoMap(regionInfoMap, regions, 2*cfg.ApiConfig.MaxInstancesToFetch)
	if err != nil {
		t.Fatal(err.Error())
	}

	cachedStartTime := time.Now()
	cachedRegionInstanceMap, err := GetInstancesRegionInfoMap(
		&cfg.ApiConfig,
		regions,
		&cfg.Credentials.Test,
		cache,
		logger,
	)
	if err != nil {
		t.Fatalf("Error thrown when fecthing instances for a second time: %s", err.Error())
	}
	cachedEndTime := time.Now()

	validateCachedRegionInfoMap(
		regionInfoMap,
		cachedRegionInstanceMap,
		regions,
		2*cfg.ApiConfig.MaxInstancesToFetch,
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

func validateRegionInfoMap(
	regionInfoMap instPkg.RegionInfoMap,
	regions []types.Region,
	maxTotalInstances int,
) error {

	for _, region := range regions {

		info, ok := regionInfoMap[region]

		if !ok || len(info.AllInstances.Instances) == 0 {
			return fmt.Errorf("Zero instances returned for region %s", REGION1)
		}

		if len(info.AllInstances.Instances) == 0 {
			return fmt.Errorf("Zero length slice for AllInstances.Instances in info for region %s", REGION1)
		}
		if len(info.PermanentInstances.Instances) == 0 {
			return fmt.Errorf("Zero length slice for PermanentInstances.Instances in info for region %s", REGION1)
		}

		if len(info.AllInstances.Instances) != info.AllInstances.Aggregates.Count {
			return fmt.Errorf(
				"Aggregates count does not match number of instances for AllInstances in region %s. "+
					"Number of instances: %d, Aggregates count: %d",
				region.ToCodeString(),
				len(info.AllInstances.Instances),
				info.AllInstances.Aggregates.Count,
			)
		}
		if len(info.PermanentInstances.Instances) != info.PermanentInstances.Aggregates.Count {
			return fmt.Errorf(
				"Aggregates count does not match number of instances for PermanentInstances in region %s. "+
					"Number of instances: %d, Aggregates count: %d",
				region.ToCodeString(),
				len(info.PermanentInstances.Instances),
				info.PermanentInstances.Aggregates.Count,
			)
		}

		if len(info.AllInstances.Instances) > maxTotalInstances {
			return fmt.Errorf(
				"More instances returned than config max for region %s. Wanted: < %d, got: %d",
				region.ToCodeString(),
				maxTotalInstances,
				len(info.AllInstances.Instances),
			)
		}
	}

	return nil
}

func validateCachedRegionInfoMap(
	fetchedRegionInfoMap instPkg.RegionInfoMap,
	storedRegionInfoMap instPkg.RegionInfoMap,
	regions []types.Region,
	maxTotalInstances int,
) error {

	validateRegionInfoMap(fetchedRegionInfoMap, regions, maxTotalInstances)

	for _, region := range regions {
		storedInfo := storedRegionInfoMap[region]
		fetchedInfo := fetchedRegionInfoMap[region]

		if len(storedInfo.AllInstances.Instances) != len(fetchedInfo.AllInstances.Instances) {
			return fmt.Errorf(
				"Different number of instances (AllInstances) fetched from cache than stored. "+
					"Region: %s, stored: %d, fetched: %d",
				region.ToCodeString(),
				len(storedInfo.AllInstances.Instances),
				len(fetchedInfo.AllInstances.Instances),
			)
		}

		if len(storedInfo.PermanentInstances.Instances) != len(fetchedInfo.PermanentInstances.Instances) {
			return fmt.Errorf(
				"Different number of instances (PermanentInstances) fetched from cache than stored. "+
					"Region: %s, stored: %d, fetched: %d",
				region.ToCodeString(),
				len(storedInfo.PermanentInstances.Instances),
				len(fetchedInfo.PermanentInstances.Instances),
			)
		}

		if storedInfo.AllInstances.Aggregates.Count != fetchedInfo.PermanentInstances.Aggregates.Count {
			return fmt.Errorf(
				"Different aggregate values fetched from cache than stored. Region: %s, stored: %+v, fetched:%+v",
				region.ToCodeString(),
				storedInfo.AllInstances.Aggregates.Count,
				fetchedInfo.PermanentInstances.Aggregates.Count,
			)
		}
	}

	return nil
}
