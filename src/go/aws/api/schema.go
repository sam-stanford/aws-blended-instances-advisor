package api

// TODO: Rename these
import (
	. "ec2-test/aws/types"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
)

type onDemandInstanceInfo struct {
	Specs       onDemandInstanceSpecs   `json:"product"`
	ServiceCode string                  `json:"serviceCode"`
	Pricing     onDemandInstancePricing `json:"terms"`
}

type onDemandInstanceSpecs struct {
	Family     string                     `json:"productFamily"`
	Sku        string                     `json:"sku"`
	Attributes onDemandInstanceAttributes `json:"attributes"`
}

type onDemandInstanceAttributes struct {
	AvailabilityZone           string `json:"availabilityzone"`
	CapacityStatus             string `json:"capacitystatus"`
	ClassicNetworkingSupport   string `json:"classicnetworkingsupport"`
	ClockSpeed                 string `json:"clockSpeed"`
	CurrentGeneration          string `json:"currentGeneration"`
	DedicatedEbsThroughput     string `json:"dedicatedEbsThroughput"`
	Ecu                        string `json:"ecu"`
	EnhancedNetorkingSupported string `json:"enhancedNetworkingSupported"`
	InstanceFamily             string `json:"instanceFamily"`
	InstanceSku                string `json:"instancesku"`
	InstanceType               string `json:"instanceType"`
	IntelAvx2Available         string `json:"intelAvx2Available"`
	IntelAvxAvailable          string `json:"intelAvxAvailable"`
	IntelTurboAvailable        string `json:"intelTurboAvailable"`
	LicenseModel               string `json:"licenseModel"`
	Location                   string `json:"location"`
	LocationType               string `json:"locationType"`
	MarketOption               string `json:"marketoption"`
	Memory                     string `json:"memory"`
	NetworkPerformance         string `json:"networkPerformance"`
	NormalizationSizeFactor    string `json:"normalizationSizeFactor"`
	OperatingSystem            string `json:"operatingSystem"`
	Operation                  string `json:"operation"`
	PhysicalProcessor          string `json:"physicalProcessor"`
	PreInstalledSw             string `json:"preInstalledSw"`
	ProcessorArchitecture      string `json:"processorArchitecture"`
	ProcessorFeatures          string `json:"processorFeatures"`
	ServiceCode                string `json:"servicecode"`
	ServiceName                string `json:"servicename"`
	Storage                    string `json:"storage"`
	Tenancy                    string `json:"tenancy"`
	UsageModel                 string `json:"usageModel"`
	Vcpu                       string `json:"vcpu"`
	VpcNetworkingSupport       string `json:"vpcnetworkingsupport"`
}

type onDemandInstancePricing struct {
	Prices map[string]onDemandInstancePricingOption `json:"OnDemand"`
}

type onDemandInstancePricingOption struct {
	Options       map[string]onDemandInstancePricingOptionDetails `json:"priceDimensions"`
	Sku           string                                          `json:"sku"`
	EffectiveDate string                                          `json:"effectiveDate"`
	OfferTermCode string                                          `json:"offerTermCode"`
}

type onDemandInstancePricingOptionDetails struct {
	Unit         string `json:"unit"`
	Description  string `json:"description"`
	RateCode     string `json:"rateCode"`
	BeginRange   string `json:"beginRange"`
	EndRange     string `json:"endRange"`
	PricePerUnit price  `json:"pricePerUnit"`
}

type price struct {
	USD string `json:"USD"`
}

type spotInstancesInfo struct {
	SpecsMap     map[string]spotInstanceSpecs                `json:"instance_types"`
	RegionPrices map[string]regionSpotInstanceRevocationInfo `json:"spot_advisor"`
}

type spotInstanceSpecs struct {
	MemoryGb float32 `json:"ram_gb"`
	Vcpus    int     `json:"cores"`
	Emr      bool    `json:"emr"`
}

type regionSpotInstanceRevocationInfo struct {
	LinuxInstances   map[string]spotInstanceRevocationInfo `json:"Linux"`
	WindowsInstances map[string]spotInstanceRevocationInfo `json:"Windows"`
}

type spotInstanceRevocationInfo struct {
	// TODO: Move comments to doc
	RevocationProbabilityTier int `json:"r"` // 0 => <5%, 1 => 5-10%, 2 => 10-15%, 3 => 15-20%, 4 => >20%
	PercentageSavings         int `json:"s"` // Over on-demand
}

func (info *onDemandInstanceInfo) toInstance() (*Instance, error) {
	vcpus, err := parseOnDemandVcpus(info)
	if err != nil {
		return nil, err
	}
	mem, err := parseOnDemandMemory(info)
	if err != nil {
		return nil, err
	}
	region, err := NewRegionFromString(info.Specs.Attributes.Location)
	if err != nil {
		return nil, err
	}
	os, err := NewOperatingSystemFromString(info.Specs.Attributes.OperatingSystem)
	if err != nil {
		return nil, err
	}
	price, err := parseOnDemandPrice(info)
	if err != nil {
		return nil, err
	}

	return &Instance{
		Name:                  info.Specs.Attributes.InstanceType,
		MemoryGb:              mem,
		Vcpus:                 vcpus,
		Region:                region,
		AvailabilityZone:      info.Specs.Attributes.AvailabilityZone,
		OperatingSystem:       os,
		PricePerHour:          price,
		RevocationProbability: 0, // On-demand instances will not be revoked
	}, nil
}

func parseOnDemandApiResponseToInstances(resp *pricing.GetProductsOutput) ([]Instance, error) {

	instances := make([]Instance, 0)

	for _, instanceInfoJson := range resp.PriceList {
		var info onDemandInstanceInfo
		err := json.Unmarshal([]byte(instanceInfoJson), &info)
		if err != nil {
			return nil, err
		}

		if info.Specs.Attributes.MarketOption == "OnDemand" {
			instance, err := info.toInstance()
			if err != nil {
				return nil, err // TODO: Handle this more gracefully
			}

			instances = append(instances, *instance)
		}
	}

	return instances, nil
}

func parseOnDemandVcpus(info *onDemandInstanceInfo) (int, error) {
	return strconv.Atoi(info.Specs.Attributes.Vcpu)
}

func parseOnDemandMemory(info *onDemandInstanceInfo) (float32, error) {
	// TODO: Manage non-GB / non-GiB values
	memStr := info.Specs.Attributes.Memory

	index := 0
	for index < len(memStr) && isNumber(memStr[index]) {
		index += 1
	}

	mem, err := strconv.ParseFloat(memStr[:index], 32)
	return float32(mem), err
}

func isNumber(b byte) bool {
	return b >= '0' && b <= '9'
}

func parseOnDemandPrice(info *onDemandInstanceInfo) (float64, error) {
	errStr := "more than one price provided when parsing on-demand instance price"

	if len(info.Pricing.Prices) != 1 {
		return -1, errors.New(errStr)
	}

	for _, price := range info.Pricing.Prices {
		if len(price.Options) > 1 {
			return -1, errors.New(errStr)
		}
		for _, option := range price.Options {
			return strconv.ParseFloat(option.PricePerUnit.USD, 64)
		}
	}
	return -1, nil
}

func parseSpotPrice(info *ec2Types.SpotPrice) (float64, error) {
	return strconv.ParseFloat(*info.SpotPrice, 64)
}

func (info *spotInstanceRevocationInfo) getRevocationProbability() (float32, error) {
	// Return the upper bound of the tier
	switch info.RevocationProbabilityTier {
	case 0:
		return 0.05, nil
	case 1:
		return 0.1, nil
	case 2:
		return 0.15, nil
	case 3:
		return 0.2, nil
	case 4:
		return 0.3, nil // TODO: >20% => ?
	default:
		return -1, errors.New(
			fmt.Sprintf(
				"provided revocation probability tier does not exist: %d",
				info.RevocationProbabilityTier,
			))
	}
}
