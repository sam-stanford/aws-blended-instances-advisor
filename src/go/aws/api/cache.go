package api

import (
	"ec2-test/cache"
	instPkg "ec2-test/instances"
	"encoding/json"
	"errors"
	"time"
)

func getRegionInfoMapFromCache(instancesCacheFilename string, c *cache.Cache) (instPkg.RegionInfoMap, error) {
	isValid := c.IsValid(instancesCacheFilename)
	if isValid {
		instancesFileContent, err := c.Get(instancesCacheFilename)
		if err != nil {
			return nil, err
		}
		var regionInfoMap instPkg.RegionInfoMap
		err = json.Unmarshal([]byte(instancesFileContent), &regionInfoMap)
		if err != nil {
			return nil, err
		}
		return regionInfoMap, nil
	}
	return nil, errors.New("instances not in cache")
}

func storeRegionInfoMapInCache(
	regionInfoMap instPkg.RegionInfoMap,
	instancesCacheFilename string,
	c *cache.Cache,
) error {
	instancesFileContent, err := json.Marshal(regionInfoMap)
	if err != nil {
		return err
	}
	err = c.Set(
		instancesCacheFilename,
		string(instancesFileContent),
		time.Hour*INSTANCES_CACHE_DURATION,
	)
	if err != nil {
		return err
	}
	return nil
}
