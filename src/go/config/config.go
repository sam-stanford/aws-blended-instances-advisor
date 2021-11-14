package config

import (
	awsTypes "ec2-test/aws/types"
	"ec2-test/utils"
	"encoding/json"
)

type Credentials struct {
	AwsKeyId     string `json:"awsKeyId"`
	AwsSecretKey string `json:"awsSecretKey"`
}

type Service struct {
	Name                  string  `json:"name"`
	MinMemory             float64 `json:"minMemory"`
	MinVcpu               int     `json:"minVcpu"`
	RevocationSensitivity float64 `json:"revocationSensitivity"`
}

type Constraints struct {
	Regions          []string  `json:"regions"`
	Services         []Service `json:"services"`
	MaxInstanceCount int       `json:"maxInstanceCount"`
}

type Endpoints struct {
	AwsSpotInstanceInfoUrl string `json:"awsSpotInstanceInfoUrl"`
}

type Config struct {
	Credentials  Credentials `json:"credentials"`
	Constraints  Constraints `json:"constraints"`
	Endpoints    Endpoints   `json:"endpoints"`
	DownloadsDir string      `json:"downloadsDir"`
}

type publicConfig struct {
	Constraints  Constraints `json:"constraints"`
	Endpoints    Endpoints   `json:"endpoints"`
	DownloadsDir string      `json:"downloadsDir"`
}

// TODO: Doc OR improve
func (c *Config) ToPublicJson() string {
	publicConfig := publicConfig{
		Constraints:  c.Constraints,
		Endpoints:    c.Endpoints,
		DownloadsDir: c.DownloadsDir,
	}
	jsonBytes, _ := json.Marshal(publicConfig)
	return string(jsonBytes)
}

func GetConfig(filepath string) (*Config, error) {
	configBytes, err := utils.FileToBytes(filepath)
	if err != nil {
		return GetEmptyConfig(), err
	}

	var c Config
	err = json.Unmarshal(configBytes, &c)
	if err != nil {
		return GetEmptyConfig(), err
	}
	return &c, nil
}

func GetEmptyConfig() *Config {
	return &Config{}
}

func (config *Config) GetRegions() []awsTypes.Region {
	regions := make([]awsTypes.Region, len(config.Constraints.Regions))
	for _, regionStr := range config.Constraints.Regions {
		region, err := awsTypes.NewRegionFromString(regionStr)
		if err == nil {
			// TODO: Do something with this error?
			regions = append(regions, region)
		}
	}
	return regions
}
