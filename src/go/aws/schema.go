package aws

type productAttributes struct {
	AvailabilityZone           string `json:"availabilityzone"`
	CapacityStatus             string `json:"capacitystatus"`
	ClassicNetworkingSupport   string `json:"classicnetworkingsupport"`
	ClockSpeed                 string `json:"clockSpeed"`
	CurrentGeneration          string `json:"currentGeneration"`
	DedicatedEbsThroughput     string `json:"dedicatedEbsThroughput"`
	Ecu                        string `json:"ecu"` // Float in quotes
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
	NormalizationSizeFactor    string `json:"normalizationSizeFactor"` // Float in quotes
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
	Vcpu                       string `json:"vcpu"`                 // Int in quotes
	VpcNetworkingSupport       string `json:"vpcnetworkingsupport"` // Bool in quotes
}

type productInfo struct {
	Family     string            `json:"productFamily"`
	Sku        string            `json:"sku"`
	Attributes productAttributes `json:"attributes"`
}

type price struct {
	USD string `json:"USD"`
}

type productPricingOption struct {
	Unit         string `json:"unit"`
	Description  string `json:"description"`
	RateCode     string `json:"rateCode"`
	BeginRange   string `json:"beginRange"`
	EndRange     string `json:"endRange"`
	PricePerUnit price  `json:"pricePerUnit"`
}

type productPricingInfo struct {
	Options       map[string]productPricingOption `json:"priceDimensions"`
	Sku           string                          `json:"sku"`
	EffectiveDate string                          `json:"effectiveDate"`
	OfferTermCode string                          `json:"offerTermCode"`
}

type product struct {
	Info        productInfo                              `json:"product"`
	ServiceCode string                                   `json:"serviceCode"`
	Pricing     map[string]map[string]productPricingInfo `json:"terms"` // E.g. "OnDemand" > "22238NQZQCYQYQ6B.JRTCKXETXF" > PricingInfo
}
