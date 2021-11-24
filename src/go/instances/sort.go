package instances

import (
	"sort"
	"strings"
)

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

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of their revocation probabilities.
func SortInstancesByRevocationProbability(instances []Instance, startIndex int, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[i].RevocationProbability < instances[j].RevocationProbability
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing lexographcial order of their operating system.
func SortInstancesByOperatingSystem(instances []Instance, startIndex int, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return strings.Compare(
			instances[i].OperatingSystem.ToString(),
			instances[j].OperatingSystem.ToString(),
		) == -1
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing lexographcial order of their Region code name.
func SortInstancesByRegion(instances []Instance, startIndex int, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return strings.Compare(
			instances[i].Region.ToCodeString(),
			instances[j].Region.ToCodeString(),
		) == -1
	})
}
