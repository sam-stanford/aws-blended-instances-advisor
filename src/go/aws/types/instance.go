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

/* TODO: Notes

1) For each region in config...
2) Considering only instances with >= mem than config...
3) Ingoring apps that must have at least on on-demand for now...
4) Suggest instance that has the best computation-revocation-price combination
	- Computation & revocation => Time preference
		- Can build relationship between computation & revocation
	- Price => Monetary preference

Time preference
- Naive: Take most powerful on-demand instance
	- BUT we can achieve same performance for cheaper
		- This is the key challenge

*/
