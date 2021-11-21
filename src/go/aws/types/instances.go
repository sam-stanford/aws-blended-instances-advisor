package types

import "sort"

// TODO: Doc
type Instance struct {
	Name                  string          `json:"name"`
	MemoryGb              float32         `json:"memory"`
	Vcpus                 int             `json:"vcpu"`
	Region                Region          `json:"region"`
	AvailabilityZone      string          `json:"az"`
	OperatingSystem       OperatingSystem `json:"os"`
	PricePerHour          float64         `json:"price"`
	RevocationProbability float32         `json:"revocProb"`
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of price.
func SortInstancesByPrice(instances []Instance, startIndex int, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[i].PricePerHour < instances[j].PricePerHour
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of memory.
func SortInstancesByMemory(instances []Instance, startIndex int, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[i].MemoryGb < instances[j].MemoryGb
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of the number of VCPUs.
func SortInstancesByVcpus(instances []Instance, startIndex int, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[i].Vcpus < instances[j].Vcpus
	})
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
