package instances

import (
	"sort"
	"strings"
)

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of price.
func SortInstancesByPrice(instances []Instance, startIndex int, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].PricePerHour < instances[startIndex+j].PricePerHour
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of memory.
func SortInstancesByMemory(instances []Instance, startIndex int, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].MemoryGb < instances[startIndex+j].MemoryGb
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of the number of VCPUs.
func SortInstancesByVcpus(instances []Instance, startIndex int, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].Vcpus < instances[startIndex+j].Vcpus
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of their revocation probabilities.
func SortInstancesByRevocationProbability(instances []Instance, startIndex int, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].RevocationProbability < instances[startIndex+j].RevocationProbability
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing lexographcial order of their operating system.
func SortInstancesByOperatingSystem(instances []Instance, startIndex int, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return strings.Compare(
			instances[startIndex+i].OperatingSystem.ToString(),
			instances[startIndex+j].OperatingSystem.ToString(),
		) == -1
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing lexographcial order of their Region code name.
func SortInstancesByRegion(instances []Instance, startIndex int, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return strings.Compare(
			instances[startIndex+i].Region.ToCodeString(),
			instances[startIndex+j].Region.ToCodeString(),
		) == -1
	})
}
