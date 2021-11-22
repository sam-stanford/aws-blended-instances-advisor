package config

import (
	"math"
	"strings"
	"testing"
)

const (
	TEST_CONFIG_FILEPATH = "../../../assets/test/example-config.json"
)

func parseConfigAndCheckForErrors(t *testing.T) *Config {
	cfg, err := ParseConfig(TEST_CONFIG_FILEPATH)
	if err != nil {
		t.Fatalf("Error thrown when parsing test config: %s", err.Error())
	}
	if cfg == nil {
		t.Fatalf("Nil config struct returned when parsing test config")
	}
	return cfg
}

func TestParseConfig(t *testing.T) {
	cfg := parseConfigAndCheckForErrors(t)

	if cfg.Credentials.Development.AwsKeyId != "DEV_KEY_ID" {
		t.Fatalf("Incorrect AWS key ID parsed. Wanted: %s, got: %s", "DEV_KEY_ID", cfg.Credentials.Development.AwsKeyId)
	}
	if cfg.Credentials.Development.AwsSecretKey != "DEV_SECRET_KEY" {
		t.Fatalf("Incorrect AWS secret key parsed. Wanted: %s, got: %s", "DEV_SECRET_KEY", cfg.Credentials.Development.AwsSecretKey)
	}
	if cfg.Credentials.Production.AwsKeyId != "PROD_KEY_ID" {
		t.Fatalf("Incorrect AWS key ID parsed. Wanted: %s, got: %s", "PROD_KEY_ID", cfg.Credentials.Production.AwsKeyId)
	}
	if cfg.Credentials.Production.AwsSecretKey != "PROD_SECRET_KEY" {
		t.Fatalf("Incorrect AWS secret key parsed. Wanted: %s, got: %s", "PROD_SECRET_KEY", cfg.Credentials.Production.AwsSecretKey)
	}
	if cfg.Credentials.Test.AwsKeyId != "TEST_KEY_ID" {
		t.Fatalf("Incorrect AWS key ID parsed. Wanted: %s, got: %s", "TEST_KEY_ID", cfg.Credentials.Test.AwsKeyId)
	}
	if cfg.Credentials.Test.AwsSecretKey != "TEST_SECRET_KEY" {
		t.Fatalf("Incorrect AWS secret key parsed. Wanted: %s, got: %s", "TEST_SECRET_KEY", cfg.Credentials.Test.AwsSecretKey)
	}

	if len(cfg.Constraints.Regions) != 2 {
		t.Fatalf("Incorrect number of regions parsed. Wanted: %d, got: %d", 2, len(cfg.Constraints.Regions))
	}
	if cfg.Constraints.Regions[0] != "us-west-1" {
		t.Fatalf("Incorrect region parsed. Wanted: %s, got: %s", "us-west-1", cfg.Constraints.Regions[0])
	}
	if cfg.Constraints.Regions[1] != "us-west-2" {
		t.Fatalf("Incorrect region parsed. Wanted: %s, got: %s", "us-west-2", cfg.Constraints.Regions[1])
	}

	if len(cfg.Constraints.Services) != 1 {
		t.Fatalf("Incorrect number of services parsed. Wanted: %d, got: %d", 1, len(cfg.Constraints.Services))
	}
	service := cfg.Constraints.Services[0]
	if service.Name != "TS1" {
		t.Fatalf("Incorrect service name parsed. Wanted: %s, got: %s", "TS1", service.Name)
	}
	if math.Abs(service.MinMemory-2.0) > 0.001 {
		t.Fatalf("Incorrect service memory parsed. Wanted: %f, got: %f", 2.0, service.MinMemory)
	}
	if service.MinVcpu != 1 {
		t.Fatalf("Incorrect service VCPU parsed. Wanted: %d, got: %d", 1, service.MinVcpu)
	}
	if math.Abs(service.RevocationSensitivity-0.1) > 0.001 {
		t.Fatalf("Incorrect service revocation sensitivity parsed. Wanted: %f, got: %f", 0.1, service.RevocationSensitivity)
	}

	if cfg.ApiConfig.DownloadsDir != "TEST_DIR" {
		t.Fatalf("Incorrect downloads directory parsed. Wanted: %s, got: %s", "TEST_DIR", cfg.ApiConfig.DownloadsDir)
	}
	if cfg.ApiConfig.MaxInstancesToFetch != 1000 {
		t.Fatalf("Incorrect max instances to fetch value parsed. Wanted: %d, got: %d", 1000, cfg.ApiConfig.MaxInstancesToFetch)
	}
	if cfg.ApiConfig.Endpoints.AwsSpotInstanceInfoUrl != "TEST_URL" {
		t.Fatalf("Incorrect spot instances URL parsed. Wanted: %s, got: %s", "TEST_URL", cfg.ApiConfig.Endpoints.AwsSpotInstanceInfoUrl)
	}

	if cfg.CacheConfig.Dirpath != "TEST_CACHE_DIRPATH" {
		t.Fatalf("Incorrect cache filepath parsed. Wanted: %s, got: %s", "TEST_CACHE_DIRPATH", cfg.CacheConfig.Dirpath)
	}
	if cfg.CacheConfig.DefaultLifetime != 200 {
		t.Fatalf("Incorrect cache default lifetime parsed. Wanted: %d, got: %d", 200, cfg.CacheConfig.DefaultLifetime)
	}
}

func TestToStringDoesNotContainCredentials(t *testing.T) {
	keyId, secretKey := "KEY_ID", "SECRET_KEY"
	cfg := Config{
		Credentials: CredentialsList{
			Production: Credentials{
				AwsKeyId:     keyId,
				AwsSecretKey: secretKey,
			},
			Development: Credentials{
				AwsKeyId:     keyId,
				AwsSecretKey: secretKey,
			},
			Test: Credentials{
				AwsKeyId:     keyId,
				AwsSecretKey: secretKey,
			},
		},
	}
	cfgJson := cfg.ToString()
	if strings.Contains(cfgJson, keyId) {
		t.Fatalf("JSON string contains the AWS key ID: %s", cfgJson)
	}
	if strings.Contains(cfgJson, secretKey) {
		t.Fatalf("JSON string contains the AWS secret key: %s", cfgJson)
	}
}

func TestValidationErrorForMissingCredentials(t *testing.T) {
	cfg := parseConfigAndCheckForErrors(t)
	cfg.Credentials.Development.AwsKeyId = ""

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("no error thrown for missing credentials during config validation")
	}
}

func TestValidationErrorForMissingApiConfig(t *testing.T) {
	cfg := parseConfigAndCheckForErrors(t)
	cfg.ApiConfig.Endpoints.AwsSpotInstanceInfoUrl = ""

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("no error thrown for missing API endpoint during config validation")
	}
}

func TestErorForInvalidRegions(t *testing.T) {
	cfg := parseConfigAndCheckForErrors(t)
	cfg.Constraints.Regions = []string{"INVALID REGION"}

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("no error thrown for invalid region during config validation")
	}
}

func TestValidationErrorForNoRegions(t *testing.T) {
	cfg := parseConfigAndCheckForErrors(t)
	cfg.Constraints.Regions = make([]string, 0)

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("no error thrown for no regions during config validation")
	}
}

func TestValidationErrorForNoServiceDescriptions(t *testing.T) {
	cfg := parseConfigAndCheckForErrors(t)
	cfg.Constraints.Services = make([]ServiceDescription, 0)

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("no error thrown for no service descriptions during config validation")
	}
}

func TestValidationErrorForNoServiceName(t *testing.T) {
	cfg := parseConfigAndCheckForErrors(t)
	cfg.Constraints.Services[0].Name = ""

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("no error thrown for empty service name during config validation")
	}
}

func TestValidationErrorForDuplicateServiceNames(t *testing.T) {
	cfg := parseConfigAndCheckForErrors(t)
	cfg.Constraints.Services = []ServiceDescription{
		{
			Name:                  "DUPLICATE NAME",
			MinMemory:             0,
			MinVcpu:               1,
			RevocationSensitivity: 0.5,
		},
		{
			Name:                  "DUPLICATE NAME",
			MinMemory:             0,
			MinVcpu:               1,
			RevocationSensitivity: 0.5,
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("no error thrown for duplicate service name during config validation")
	}
}

func TestValidationErrorForInvalidServiceRevocationSensitivity(t *testing.T) {
	cfg := parseConfigAndCheckForErrors(t)
	cfg.Constraints.Services[0].RevocationSensitivity = -50

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("no error thrown for invalid service revocation sensitivity during config validation")
	}

	cfg.Constraints.Services[0].RevocationSensitivity = 1.01

	err = cfg.Validate()
	if err == nil {
		t.Fatalf("no error thrown for invalid service revocation sensitivity during config validation")
	}
}
