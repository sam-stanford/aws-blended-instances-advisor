package config

import (
	awsTypes "ec2-test/aws/types"
	"ec2-test/utils"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	DEFAULT_API_PORT                       = 12021
	DEFAULT_AWS_API_SPOT_INSTANCE_INFO_URL = "https://spot-bid-advisor.s3.amazonaws.com/spot-advisor-data.json"
	DEFAULT_AWS_API_MAX_INSTANCES_TO_FETCH = 0
	DEFAULT_AWS_API_DOWNLOADS_DIR          = "../../assets/downloads"
	DEFAULT_CACHE_DIR                      = "../../assets/cache"
	DEFAULT_CACHE_DEFAULT_LIFETIME         = 96
)

// TODO: Doc

// TODO: Update actual, example, and test configs

// TODO: Remove servicess

type Config struct {
	ApiConfig    ApiConfig       `json:"api"`
	AwsApiConfig AwsApiConfig    `json:"awsApi"`
	CacheConfig  CacheConfig     `json:"cache"`
	Constraints  Constraints     `json:"constraints"` // TODO: Remove this (add max instances to AwsApiConfig)
	Credentials  CredentialsList `json:"credentials"`
}

type ApiConfig struct {
	Port int `json:"port"`
}

type AwsApiConfig struct { // TODO: Potentially rename to generic config
	Endpoints             Endpoints `json:"endpoints"`
	DownloadsDir          string    `json:"downloadsDir"`
	MaxInstancesToFetch   int       `json:"maxInstancesToFetch"`
	ConsiderFreeInstances bool      `json:"freeInstances"`
}

type Endpoints struct {
	AwsSpotInstanceInfoUrl string `json:"awsSpotInstanceInfoUrl"`
}

type CacheConfig struct {
	Dirpath string `json:"dirpath"`
}

type CredentialsList struct {
	Production  Credentials `json:"production"`
	Development Credentials `json:"development"`
	Test        Credentials `json:"test"`
}

type Credentials struct {
	AwsKeyId     string `json:"awsKeyId"`
	AwsSecretKey string `json:"awsSecretKey"`
}

type Constraints struct {
	Regions  []string             `json:"regions"` // TODO: Move to config
	Services []ServiceDescription `json:"services"`
}

type ServiceDescription struct {
	Name        string                      `json:"name"`
	MinMemory   float64                     `json:"minMemory"`
	MaxVcpu     int                         `json:"maxVcpu"`
	Instances   ServiceDescriptionInstances `json:"instances"`
	Focus       string                      `json:"focus"`
	FocusWeight float64                     `json:"focusWeight"`
}

type ServiceDescriptionInstances struct {
	MinimumCount int `json:"minimum"`
	TotalCount   int `json:"total"`
}

// TODO: Doc comments

func (c *Constraints) GetRegions() []awsTypes.Region {
	regions := []awsTypes.Region{}
	for _, regionStr := range c.Regions {
		region, err := awsTypes.NewRegion(regionStr)
		if err == nil {
			regions = append(regions, region)
		}
	}
	return regions
}

func (c *Config) String() string { // TODO: Use in main logging
	noCredsConfig := &Config{
		Constraints: c.Constraints,
		ApiConfig:   c.ApiConfig,
		CacheConfig: c.CacheConfig,
		Credentials: c.Credentials,
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
			Dirpath: DEFAULT_CACHE_DIR,
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

	err := c.validateConstraints()
	if err != nil {
		return err
	}
	err = c.validateCredentials()
	if err != nil {
		return err
	}
	err = c.validateApiConfig()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) validateConstraints() error {
	err := c.validateRegions()
	if err != nil {
		return err
	}
	err = c.validateServiceDescriptions()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) validateRegions() error {
	regions := c.Constraints.Regions
	if len(regions) == 0 {
		return errors.New("invalid config: zero regions provided")
	}
	for _, r := range regions {
		_, err := awsTypes.NewRegion(r)
		if err != nil {
			return utils.PrependToError(err, "invalid config")
		}
	}
	return nil
}

func (c *Config) validateServiceDescriptions() error {
	services := c.Constraints.Services
	if len(services) == 0 {
		return errors.New("invalid config: zero service descriptions provided")
	}
	dupeName := findFirstDuplicateServiceName(services)
	if dupeName != "" {
		return fmt.Errorf("invalid config: duplicate service name: %s", dupeName)
	}

	for _, s := range services {
		err := validateServiceDescription(s)
		if err != nil {
			return utils.PrependToError(err, "invalid config")
		}
	}

	return nil
}

func validateServiceDescription(svc ServiceDescription) error {
	if svc.Name == "" {
		return errors.New("service name is empty")
	}
	return nil
}

func findFirstDuplicateServiceName(services []ServiceDescription) string {
	names := make(map[string]bool)

	for _, s := range services {
		_, exists := names[s.Name]
		if exists {
			return s.Name
		}
		names[s.Name] = true
	}

	return ""
}

func (c *Config) validateCredentials() error {
	// No validation required
	return nil
}

func (c *Config) validateApiConfig() error {
	// No validation required
	return nil
}
