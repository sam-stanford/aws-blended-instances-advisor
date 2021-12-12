//go:build integration
// +build integration

package api

import (
	"ec2-test/aws/types"
	"ec2-test/config"
	"ec2-test/utils"
	"testing"
)

func TestGetOnDemandInstances(t *testing.T) {
	cfg, err := config.ParseConfig(CONFIG_FILEPATH)
	if err != nil {
		t.Fatalf("Failed to read config: %s", err.Error())
	}
	cfg.AwsApiConfig.MaxInstancesToFetch = MAX_INSTANCES
	cfg.AwsApiConfig.DownloadsDir = CONFIG_API_DOWNLOADS_DIR

	creds := createAwsCredentials(&cfg.Credentials.Test)

	logger, err := utils.CreateMockLogger()
	if err != nil {
		t.Fatalf("Failed to create mock logger: %s", err.Error())
	}

	region1, region2, err := getTestRegions()
	if err != nil {
		t.Fatalf("Error occured when parsing test regions: %s", err.Error())
	}

	regionInstanceMap, err := GetOnDemandInstances(&cfg.AwsApiConfig, []types.Region{*region1, *region2}, creds, logger)
	if err != nil {
		t.Fatalf("Error thrown when fetching on-demand instances: %s", err.Error())
	}

	instances, ok := regionInstanceMap[*region1]
	if !ok || len(instances) == 0 {
		t.Fatalf("Zero instances returned for region %s", REGION1)
	}
	if len(instances) > cfg.AwsApiConfig.MaxInstancesToFetch {
		t.Fatalf(
			"More instances returned than config max for region %s. Wanted: <%d, got: %d",
			REGION1,
			cfg.AwsApiConfig.MaxInstancesToFetch,
			len(instances),
		)
	}

	instances, ok = regionInstanceMap[*region2]
	if !ok || len(instances) == 0 {
		t.Fatalf("Zero instances returned for region %s", REGION2)
	}
	if len(instances) > cfg.AwsApiConfig.MaxInstancesToFetch {
		t.Fatalf(
			"More instances returned than config max for region %s. Wanted: <%d, got: %d",
			REGION2,
			cfg.AwsApiConfig.MaxInstancesToFetch,
			len(instances),
		)
	}
}
