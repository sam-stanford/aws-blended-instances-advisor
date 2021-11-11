package types

type Instance struct {
	Name                  string
	MemoryGb              float32
	Vcpus                 int
	Region                Region
	AvailabilityZone      string
	OperatingSystem       OperatingSystem
	PricePerHour          float64
	RevocationProbability float32
}
