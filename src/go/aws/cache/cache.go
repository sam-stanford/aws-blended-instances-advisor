package cache

import (
	"ec2-test/aws/types"
	"time"
)

func CheckInstancesCacheValid(cacheFilepath string, cacheLifetime time.Duration) bool {
	// TODO
	return false
}

func GetInstancesCache(cacheFilepath string) ([]types.Instance, error) {
	// TODO
	return nil, nil
}

func SetInstancesCache(cacheFilepath string, instances []types.Instance) error {
	// TODO
	return nil
}
