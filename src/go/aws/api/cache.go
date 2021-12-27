package schema

import (
	"aws-blended-instances-advisor/cache"
	instPkg "aws-blended-instances-advisor/instances"
	"encoding/json"
	"errors"
	"time"
)

// TODO: test

func getGlobalInstanceInfoFromCache(instancesCacheFilename string, c *cache.Cache) (*instPkg.GlobalInfo, error) {
	isValid := c.IsValid(instancesCacheFilename)
	if isValid {
		instancesFileContent, err := c.Get(instancesCacheFilename)
		if err != nil {
			return nil, err
		}
		var globalInfo instPkg.GlobalInfo
		err = json.Unmarshal([]byte(instancesFileContent), &globalInfo)
		if err != nil {
			return nil, err
		}
		return &globalInfo, nil
	}
	return nil, errors.New("instances not in cache")
}

func storeGlobalInstanceInfoInCache(
	globalInstanceInfo instPkg.GlobalInfo,
	instancesCacheFilename string,
	c *cache.Cache,
) error {
	instancesFileContent, err := json.Marshal(globalInstanceInfo)
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
