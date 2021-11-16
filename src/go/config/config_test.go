package config

import (
	"math"
	"strings"
	"testing"
)

const (
	TEST_CONFIG_FILEPATH = "../../../assets/test/test-config.json"
)

// TODO: Constraints on values (e.g. revoc sensitivity <= 1)

func TestGetConfig(t *testing.T) {
	cfg, err := GetConfig(TEST_CONFIG_FILEPATH)
	if err != nil {
		t.Fatalf("Error thrown when parsing test config: %s", err.Error())
	}
	if cfg == nil {
		t.Fatalf("Nil config struct returned when parsing test config")
	}

	if cfg.Credentials.AwsKeyId != "KEY_ID" {
		t.Fatalf("Incorrect AWS key ID parsed. Wanted: %s, got: %s", "KEY_ID", cfg.Credentials.AwsKeyId)
	}
	if cfg.Credentials.AwsSecretKey != "SECRET_KEY" {
		t.Fatalf("Incorrect AWS secret key parsed. Wanted: %s, got: %s", "SECRET_KEY", cfg.Credentials.AwsSecretKey)
	}

	if cfg.Constraints.MaxInstanceCount != 1000 {
		t.Fatalf("Incorrect max instance count parsed. Wanted: %d, got: %d", 1000, cfg.Constraints.MaxInstanceCount)
	}
	if len(cfg.Constraints.Regions) != 2 {
		t.Fatalf("Incorrect number of regions parsed. Wanted: %d, got: %d", 2, len(cfg.Constraints.Regions))
	}
	if cfg.Constraints.Regions[0] != "region1" {
		t.Fatalf("Incorrect region parsed. Wanted: %s, got: %s", "region1", cfg.Constraints.Regions[0])
	}
	if cfg.Constraints.Regions[1] != "region2" {
		t.Fatalf("Incorrect region parsed. Wanted: %s, got: %s", "region2", cfg.Constraints.Regions[1])
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

	if cfg.Endpoints.AwsSpotInstanceInfoUrl != "TEST_URL" {
		t.Fatalf("Incorrect AWS spot instance info URL parsed. Wanted: %s, got %s", "TEST_URL", cfg.Endpoints.AwsSpotInstanceInfoUrl)
	}

	if cfg.DownloadsDir != "TEST_DIR" {
		t.Fatalf("Incorrect downloads directory parsed. Wanted: %s, got: %s", "TEST_DIR", cfg.DownloadsDir)
	}
}

func TestToPublicJson(t *testing.T) {
	keyId, secretKey := "KEY_ID", "SECRET_KEY"
	cfg := Config{
		Credentials: Credentials{
			AwsKeyId:     keyId,
			AwsSecretKey: secretKey,
		},
	}
	cfgJson := cfg.ToPublicJson()
	if strings.Contains(cfgJson, keyId) {
		t.Fatalf("JSON string contains the AWS key ID: %s", cfgJson)
	}
	if strings.Contains(cfgJson, secretKey) {
		t.Fatalf("JSON string contains the AWS secret key: %s", cfgJson)
	}
}
