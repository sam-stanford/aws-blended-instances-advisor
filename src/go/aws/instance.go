package aws

type Instance struct {
	Name                  string
	MemoryGb              float32
	Vcpus                 int
	Region                Region
	AvailabilityZone      string
	OperatingSystem       string
	PricePerHour          float64
	RevocationProbability float32
}

func newOnDemandInstance(info *onDemandInstanceInfo) (*Instance, error) {
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
	err = validateOperatingSystemString(info.Specs.Attributes.OperatingSystem)
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
		OperatingSystem:       info.Specs.Attributes.OperatingSystem,
		PricePerHour:          price,
		RevocationProbability: 0, // On-demand instances will not be revoked
	}, nil
}

// func newSpotInstance() (*Instance, error) {
// return
// }
