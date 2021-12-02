package config

import (
	awsTypes "ec2-test/aws/types"
	"ec2-test/utils"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	DEFAULT_API_SPOT_INSTANCE_INFO_URL = "https://spot-bid-advisor.s3.amazonaws.com/spot-advisor-data.json"
	DEFAULT_API_MAX_INSTANCES_TO_FETCH = 0
	DEFAULT_API_DOWNLOADS_DIR          = "../../assets/downloads"
	DEFAULT_CACHE_DIR                  = "../../assets/cache"
	DEFAULT_CACHE_DEFAULT_LIFETIME     = 96
)

// TODO: Doc

type Config struct {
	Credentials CredentialsList `json:"credentials"`
	Constraints Constraints     `json:"constraints"`
	ApiConfig   ApiConfig       `json:"api"`
	CacheConfig CacheConfig     `json:"cache"`
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

type Constraints struct {
	Regions  []string             `json:"regions"`
	Services []ServiceDescription `json:"services"`
}

type ApiConfig struct {
	Endpoints           Endpoints `json:"endpoints"`
	DownloadsDir        string    `json:"downloadsDir"`
	MaxInstancesToFetch int       `json:"maxInstancesToFetch"`
}

type Endpoints struct {
	AwsSpotInstanceInfoUrl string `json:"awsSpotInstanceInfoUrl"`
}

type CacheConfig struct {
	Dirpath string `json:"dirpath"`
}

// TODO: Doc comments

func (c *Constraints) GetRegions() []awsTypes.Region {
	regions := make([]awsTypes.Region, len(c.Regions))
	for _, regionStr := range c.Regions {
		region, err := awsTypes.NewRegion(regionStr)
		if err == nil {
			regions = append(regions, region)
		}
	}
	return regions
}

func (svc *ServiceDescription) GetFocus() ServiceFocus {
	return ServiceFocusFromString(svc.Focus)
}

func (c *Config) String() string {
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
			Endpoints: Endpoints{
				AwsSpotInstanceInfoUrl: DEFAULT_API_SPOT_INSTANCE_INFO_URL,
			},
			DownloadsDir:        DEFAULT_API_DOWNLOADS_DIR,
			MaxInstancesToFetch: DEFAULT_API_MAX_INSTANCES_TO_FETCH,
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
	err := ValidateServiceFocus(svc.Focus)
	if err != nil {
		return err
	}
	if svc.FocusWeight < 0 || svc.FocusWeight > 1 {
		return fmt.Errorf(
			"svc (%s) has focusWeight value outside of range of [0,1]: %f",
			svc.Name,
			svc.FocusWeight,
		)
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
