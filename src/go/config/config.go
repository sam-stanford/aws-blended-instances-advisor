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
	DEFAULT_API_CACHE_DIR              = "../../assets/cache"
)

type Config struct {
	Credentials CredentialsList `json:"credentials"`
	Constraints Constraints     `json:"constraints"`
	ApiConfig   ApiConfig       `json:"api"`
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
	Name                  string  `json:"name"`
	MinMemory             float64 `json:"minMemory"`
	MinVcpu               int     `json:"minVcpu"`
	RevocationSensitivity float64 `json:"revocationSensitivity"`
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

// TODO: Doc comments

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

func (config *Config) GetRegions() []awsTypes.Region {
	regions := make([]awsTypes.Region, len(config.Constraints.Regions))
	for _, regionStr := range config.Constraints.Regions {
		region, err := awsTypes.NewRegion(regionStr)
		if err == nil {
			regions = append(regions, region)
		}
	}
	return regions
}

func (c *Config) ToString() string {
	noCredsConfig := &Config{
		Constraints: c.Constraints,
		ApiConfig:   c.ApiConfig,
	}

	jsonBytes, _ := json.Marshal(noCredsConfig)
	return string(jsonBytes)
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

func validateServiceDescription(service ServiceDescription) error {
	if service.Name == "" {
		return errors.New("service name is empty")
	}
	if service.RevocationSensitivity < 0 || service.RevocationSensitivity > 1 {
		return fmt.Errorf(
			"service (%s) has revocation sensitive value outside of range of [0,1]: %f",
			service.Name,
			service.RevocationSensitivity,
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
