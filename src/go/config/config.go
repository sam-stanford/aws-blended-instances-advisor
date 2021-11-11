package config

import (
	. "ec2-test/aws/types"
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
	Regions  []string  `json:"regions"`
	Services []Service `json:"services"`
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

func (config *Config) GetRegions() []Region {
	regions := make([]Region, len(config.Constraints.Regions))
	for _, regionStr := range config.Constraints.Regions {
		region, err := NewRegionFromString(regionStr)
		if err == nil {
			// TODO: Do something with this error?
			regions = append(regions, region)
		}
	}
	return regions
}
