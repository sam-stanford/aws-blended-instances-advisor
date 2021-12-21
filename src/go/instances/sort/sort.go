package sort

import (
	. "aws-blended-instances-advisor/instances"
	"sort"
	"strings"
)

// Sorts a given slice of Instances in place from startIndex (inclusive) to endIndex (exclusive)
// in increasing order of price.
func SortInstancesByPrice(instances []*Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].PricePerHour < instances[startIndex+j].PricePerHour
	})
}

// Sorts a given slice of Instances in place from startIndex (inclusive) to endIndex (exclusive)
// in increasing order of memory.
func SortInstancesByMemory(instances []*Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].MemoryGb < instances[startIndex+j].MemoryGb
	})
}

// Sorts a given slice of Instances in place from startIndex (inclusive) to endIndex (exclusive)
// in increasing order of their VCPU.
func SortInstancesByVcpu(instances []*Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].Vcpu < instances[startIndex+j].Vcpu
	})
}

// Sorts a given slice of Instances in place from startIndex (inclusive) to endIndex (exclusive)
// in increasing order of their revocation probabilities.
func SortInstancesByRevocationProbability(instances []*Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].RevocationProbability < instances[startIndex+j].RevocationProbability
	})
}

// Sorts a given slice of Instances in place from startIndex (inclusive) to endIndex (exclusive)
// in increasing lexographcial order of their operating system.
func SortInstancesByOperatingSystem(instances []*Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return strings.Compare(
			instances[startIndex+i].OperatingSystem.ToString(),
			instances[startIndex+j].OperatingSystem.ToString(),
		) == -1
	})
}

// Sorts a given slice of Instances in place from startIndex (inclusive) to endIndex (exclusive)
// in increasing lexographcial order of their Region code name.
func SortInstancesByRegion(instances []*Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return strings.Compare(
			instances[startIndex+i].Region.CodeString(),
			instances[startIndex+j].Region.CodeString(),
		) == -1
	})
}

// TODO: Test / remove ?
// Sorts a given slice of Instances in place from startIndex (inclusive) to endIndex (exclusive)
// in increasing order of their score calculated from the provided weightings and aggregates.
func SortInstancesWeighted(
	instances []*Instance,
	aggregates Aggregates,
	startIndex,
	endIndex int,
	weights SortWeights,
) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		iScore := CalculateInstanceScoreFromWeights(
			instances[startIndex+i],
			aggregates,
			weights,
		)
		jScore := CalculateInstanceScoreFromWeights(
			instances[startIndex+j],
			aggregates,
			weights,
		)
		return iScore < jScore
	})
}

// Sorts a given slice of Instances in place from startIndex (inclusive) to endIndex (exclusive)
// in increasing order of their score calculated from the provided weightings and aggregates, with
// a limiter applied to instances' VCPUs after they exceed a maximum.
func SortInstancesWeightedWithVcpuLimiter(
	instances []*Instance,
	aggregates Aggregates,
	startIndex,
	endIndex int,
	weights SortWeights,
	maxVcpu int,
) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		iScore := CalculateInstanceScoreFromWeightsWithVcpuLimiter(
			instances[startIndex+i],
			aggregates,
			weights,
			maxVcpu,
		)
		jScore := CalculateInstanceScoreFromWeightsWithVcpuLimiter(
			instances[startIndex+j],
			aggregates,
			weights,
			maxVcpu,
		)
		return iScore < jScore
	})
}

// TODO: Doc & test / make private?
func CalculateInstanceScoreFromWeights(
	instance *Instance,
	aggregates Aggregates,
	weightings SortWeights, // TODO: Rename all "weightings" to "weights"
) float64 {
	normalisedVcpu := aggregates.NormaliseVcpu(instance.Vcpu)
	normalisedRp := aggregates.NormaliseRevocationProbability(instance.RevocationProbability)
	normalisedPrice := aggregates.NormalisePricePerHour(instance.PricePerHour)

	return weightings.VcpuWeight*normalisedVcpu +
		weightings.RevocationProbabilityWeight*normalisedRp +
		weightings.PriceWeight*normalisedPrice
}

// TODO: Doc & test
// TODO: Pointer to instance
func CalculateInstanceScoreFromWeightsWithVcpuLimiter(
	instance *Instance,
	aggregates Aggregates,
	weightings SortWeights,
	maxVcpu int,
) float64 {
	// TODO: Can we not just use modulo here?
	if instance.Vcpu >= maxVcpu {
		return calculatedInstanceScoreFromWeightsWithFixedVcpu(
			instance,
			aggregates,
			weightings,
			maxVcpu,
		)
	}
	return CalculateInstanceScoreFromWeights(
		instance,
		aggregates,
		weightings,
	)
}

func calculatedInstanceScoreFromWeightsWithFixedVcpu(
	instance *Instance,
	aggregates Aggregates,
	weights SortWeights,
	fixedVcpu int,
) float64 {
	normalisedVcpu := aggregates.NormaliseVcpu(fixedVcpu)
	normalisedRp := aggregates.NormaliseRevocationProbability(instance.RevocationProbability)
	normalisedPrice := aggregates.NormalisePricePerHour(instance.PricePerHour)

	return weights.VcpuWeight*normalisedVcpu +
		weights.RevocationProbabilityWeight*normalisedRp +
		weights.PriceWeight*normalisedPrice
}
