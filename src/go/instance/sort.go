package instance

import (
	"sort"
	"strings"
)

type SortWeightings struct {
	VcpuWeight                  float64
	RevocationProbabilityWeight float64 // TODO: Negative?
	PriceWeight                 float64
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of price.
func SortInstancesByPrice(instances []Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].PricePerHour < instances[startIndex+j].PricePerHour
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of memory.
func SortInstancesByMemory(instances []Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].MemoryGb < instances[startIndex+j].MemoryGb
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of their VCPU.
func SortInstancesByVcpu(instances []Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].Vcpu < instances[startIndex+j].Vcpu
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of their revocation probabilities.
func SortInstancesByRevocationProbability(instances []Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].RevocationProbability < instances[startIndex+j].RevocationProbability
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing lexographcial order of their operating system.
func SortInstancesByOperatingSystem(instances []Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return strings.Compare(
			instances[startIndex+i].OperatingSystem.ToString(),
			instances[startIndex+j].OperatingSystem.ToString(),
		) == -1
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing lexographcial order of their Region code name.
func SortInstancesByRegion(instances []Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return strings.Compare(
			instances[startIndex+i].Region.ToCodeString(),
			instances[startIndex+j].Region.ToCodeString(),
		) == -1
	})
}

// Sorts a given slice of Instances in place from startIndex to endIndex (exclusive)
// in increasing order of their score calculated from the provided weightings and aggregates.
func SortInstancesWeighted(
	instances []Instance,
	startIndex,
	endIndex int,
	weightings SortWeightings,
	aggregates Aggregates,
) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		iScore := CalculateInstanceScoreFromWeight(instances[startIndex+i], weightings, aggregates)
		jScore := CalculateInstanceScoreFromWeight(instances[startIndex+j], weightings, aggregates)
		return iScore < jScore
	})

}

func CalculateInstanceScoreFromWeight(
	instance Instance,
	weightings SortWeightings,
	aggregates Aggregates,
) float64 {
	normalisedVcpu := aggregates.NormaliseVcpu(instance.Vcpu)
	normalisedRp := aggregates.NormaliseRevocationProbability(instance.RevocationProbability)
	normalisedPrice := aggregates.NormalisePricePerHour(instance.PricePerHour)

	return weightings.VcpuWeight*normalisedVcpu +
		weightings.RevocationProbabilityWeight*normalisedRp +
		weightings.PriceWeight*normalisedPrice
}
