package config

import (
	"aws-blended-instances-advisor/utils"
	"encoding/json"
	"fmt"
)

const (
	DEFAULT_API_PORT                       = 12021
	DEFAULT_AWS_API_SPOT_INSTANCE_INFO_URL = "https://spot-bid-advisor.s3.amazonaws.com/spot-advisor-data.json"
	DEFAULT_AWS_API_MAX_INSTANCES_TO_FETCH = 0
	DEFAULT_AWS_API_DOWNLOADS_DIR          = "../../temp/downloads"
	DEFAULT_CACHE_DIR                      = "../../temp/cache"
	DEFAULT_CACHE_DEFAULT_LIFETIME         = 96
)

// TODO: Doc

// TODO: Update actual, example, and test configs

// TODO: Remove servicess

type Config struct {
	ApiConfig    ApiConfig    `json:"api"`
	AwsApiConfig AwsApiConfig `json:"awsApi"`
	CacheConfig  CacheConfig  `json:"cache"`
	Credentials  Credentials  `json:"credentials"`
}

type ApiConfig struct {
	Port           int      `json:"port"`
	AllowedDomains []string `json:"allowedDomains"`
}

type AwsApiConfig struct { // TODO: Potentially rename to generic config
	Endpoints           Endpoints `json:"endpoints"`
	DownloadsDir        string    `json:"downloadsDir"`
	MaxInstancesToFetch int       `json:"maxInstancesToFetch"`
}

type Endpoints struct {
	AwsSpotInstanceInfoUrl string `json:"awsSpotInstanceInfoUrl"` // TODO: Remove indirection layer
}

type CacheConfig struct {
	Dirpath         string `json:"dirpath"`
	DefaultLifetime int32  `json:"defaultLifetime"` // TODO: Use
}

type Credentials struct {
	AwsKeyId     string `json:"awsKeyId"`
	AwsSecretKey string `json:"awsSecretKey"`
}

func (c *Config) String() string { // TODO: Use in main logging
	noCredsConfig := &Config{
		ApiConfig:   c.ApiConfig,
		CacheConfig: c.CacheConfig,
	}

	jsonBytes, _ := json.Marshal(noCredsConfig)
	return string(jsonBytes)
}

func ParseConfig(filepath string) (*Config, error) {
	configBytes, err := utils.FileToBytes(filepath)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		ApiConfig: ApiConfig{
			Port: DEFAULT_API_PORT,
		},
		AwsApiConfig: AwsApiConfig{
			Endpoints: Endpoints{
				AwsSpotInstanceInfoUrl: DEFAULT_AWS_API_SPOT_INSTANCE_INFO_URL,
			},
			DownloadsDir:        DEFAULT_AWS_API_DOWNLOADS_DIR,
			MaxInstancesToFetch: DEFAULT_AWS_API_MAX_INSTANCES_TO_FETCH,
		},
		CacheConfig: CacheConfig{
			Dirpath:         DEFAULT_CACHE_DIR,
			DefaultLifetime: DEFAULT_CACHE_DEFAULT_LIFETIME,
		},
	}

	err = json.Unmarshal(configBytes, &cfg)
	if err != nil {
		return nil, err
	}
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {

	empty, fieldName := utils.AnyFieldsAreEmpty(c)
	if empty {
		return fmt.Errorf("config field is empty: %s", fieldName)
	}

	// TODO: Add more validation

	return nil
}
