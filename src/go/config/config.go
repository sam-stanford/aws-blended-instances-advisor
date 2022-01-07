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

// Config contains information on how the application should run.
//
// A Config object should be parsed using ParseConfig and should be
// injected into packages as appropriate.
type Config struct {
	ApiConfig    ApiConfig    `json:"api"`
	AwsApiConfig AwsApiConfig `json:"awsApi"`
	CacheConfig  CacheConfig  `json:"cache"`
	Credentials  Credentials  `json:"credentials"`
}

// ApiConfig contains details on how the API should be configured.
type ApiConfig struct {
	// The port on which the API should be run
	Port int `json:"port"`

	// The domains which are allowed to access the API
	AllowedDomains []string `json:"allowedDomains"`
}

// AwsApiConfig contains information on how the AWS API/SDK
// should be set up and used.
type AwsApiConfig struct {
	// The endpoints which are used by the API
	Endpoints Endpoints `json:"endpoints"`

	// The directory path where downloaded files should be saved
	DownloadsDir string `json:"downloadsDir"`

	// The maximum number of instances to fetch with each API call
	MaxInstancesToFetch int `json:"maxInstancesToFetch"`
}

// Endpoints contains the endpoints used in the AWS package.
type Endpoints struct {
	// The URL which spot instance info should be fetched fromd
	AwsSpotInstanceInfoUrl string `json:"awsSpotInstanceInfoUrl"`
}

// CacheConfig contains information for use in the Cache package.
type CacheConfig struct {
	// The directory path where cache files should be stored
	Dirpath string `json:"dirpath"`

	// The default lifetime for cache entries
	DefaultLifetime int32 `json:"defaultLifetime"`
}

// Credentials contains AWS API/SDK credentials.
type Credentials struct {
	// The key ID to be used for AWS API authentication
	AwsKeyId string `json:"awsKeyId"`

	// The secret key to be used for AWS APU authentication
	AwsSecretKey string `json:"awsSecretKey"`
}

// String converts a Config into a printable string, suitable for
// logging.
func (c *Config) String() string {
	noCredsConfig := &Config{
		ApiConfig:   c.ApiConfig,
		CacheConfig: c.CacheConfig,
	}

	jsonBytes, _ := json.Marshal(noCredsConfig)
	return string(jsonBytes)
}

// ParseConfig parses a Config from a file at the given filepath, filling in
// missing fields with defaults where applicable.
//
// Returns an error if an error is encountered when working with the filesystem,
// or critical fields are missing.
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

// Validate checks that a Config is valid, returning an error describing
// the problem if the Config is considered invalid.
func (c *Config) Validate() error {

	err := c.ApiConfig.validate()
	if err != nil {
		return utils.PrependToError(err, "API config invalid")
	}

	err = c.AwsApiConfig.validate()
	if err != nil {
		return utils.PrependToError(err, "AWS API config invald")
	}

	err = c.CacheConfig.validate()
	if err != nil {
		return utils.PrependToError(err, "cache config is invalid")
	}

	err = c.Credentials.validate()
	if err != nil {
		return utils.PrependToError(err, "credentials are invalid")
	}

	return nil
}

func (c *ApiConfig) validate() error {
	if c.Port <= 1023 {
		return fmt.Errorf("port %d is within controller assignment range", c.Port)
	}
	if c.AllowedDomains == nil {
		return fmt.Errorf(
			"allowedDomains is not specified. Use an empty array (\"[]\") for no allowed domains",
		)
	}
	return nil
}

func (c *AwsApiConfig) validate() error {
	if c.DownloadsDir == "" {
		return fmt.Errorf("downloadsDir is empty")
	}
	if c.Endpoints.AwsSpotInstanceInfoUrl == "" {
		return fmt.Errorf("awsSpotInstanceInfoUrl is empty")
	}
	return nil
}

func (c *CacheConfig) validate() error {
	if c.DefaultLifetime < 0 {
		return fmt.Errorf("defaultLifetime cannot be negative")
	}
	if c.Dirpath == "" {
		return fmt.Errorf("dirpath is empty")
	}
	return nil
}

func (c *Credentials) validate() error {
	if c.AwsKeyId == "" {
		return fmt.Errorf("awsKeyId is empty")
	}
	if c.AwsSecretKey == "" {
		return fmt.Errorf("awsSecretKey is empty")
	}
	return nil
}
