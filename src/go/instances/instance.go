package instances

import (
	awsTypes "ec2-test/aws/types"
)

type RegionInstancesMap map[awsTypes.Region][]Instance // TODO: Use

// TODO: Doc
type Instance struct {
	Name                  string                   `json:"name"`
	MemoryGb              float64                  `json:"memory"`
	Vcpus                 int                      `json:"vcpu"`
	Region                awsTypes.Region          `json:"region"`
	AvailabilityZone      string                   `json:"az"`
	OperatingSystem       awsTypes.OperatingSystem `json:"os"`
	PricePerHour          float64                  `json:"price"`
	RevocationProbability float64                  `json:"revocProb"`
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

Monetary preference
- Naive: Take the cheapest instances regardless of other characteristics
	- BUT revocations could actually increase the amount of time if the app is particularly sensitive to the info

*/
