package sort

import (
	instPkg "aws-blended-instances-advisor/instances"
	"sort"
	"strings"
)

// SortInstancesByPrice sorts a given slice of Instances in place from startIndex
// (inclusive) to endIndex (exclusive) in increasing order of price.
func SortInstancesByPrice(instances []*instPkg.Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].PricePerHour < instances[startIndex+j].PricePerHour
	})
}

// SortInstancesByMemory sorts a given slice of Instances in place from startIndex
// (inclusive) to endIndex (exclusive) in increasing order of memory.
func SortInstancesByMemory(instances []*instPkg.Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].MemoryGb < instances[startIndex+j].MemoryGb
	})
}

// SortInstancesByVcpu sorts a given slice of Instances in place from startIndex (inclusive)
// to endIndex (exclusive) in increasing order of their VCPU.
func SortInstancesByVcpu(instances []*instPkg.Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].Vcpu < instances[startIndex+j].Vcpu
	})
}

// SortInstancesByRevocationProbability sorts a given slice of Instances in place from startIndex
// (inclusive) to endIndex (exclusive) in increasing order of their revocation probabilities.
func SortInstancesByRevocationProbability(instances []*instPkg.Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return instances[startIndex+i].RevocationProbability < instances[startIndex+j].RevocationProbability
	})
}

// SortInstancesByRegion sorts a given slice of Instances in place from startIndex (inclusive)
// to endIndex (exclusive) in increasing lexographcial order of their Region code name.
func SortInstancesByRegion(instances []*instPkg.Instance, startIndex, endIndex int) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		return strings.Compare(
			instances[startIndex+i].Region.CodeString(),
			instances[startIndex+j].Region.CodeString(),
		) == -1
	})
}

// SortInstancesWeighed sorts a given slice of Instances in place from startIndex
// (inclusive) to endIndex (exclusive) in increasing order of their score
// calculated from the provided weights and aggregates.
func SortInstancesWeighted(
	instances []*instPkg.Instance,
	aggregates instPkg.Aggregates,
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

// SortInstancesWeightedWithVcpuLimiter sorts a given slice of Instances in place from
// startIndex (inclusive) to endIndex (exclusive) in increasing order of their score
// calculated from the provided weights and aggregates, with a limiter applied to
// instances' VCPUs after they exceed a maximum.
func SortInstancesWeightedWithVcpuLimiter(
	instances []*instPkg.Instance,
	aggregates instPkg.Aggregates,
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

// CalculateInstanceScoreFromWeights computes a score for an Instance relative
// to other Instances by using a given set of SortWeights and Instance Aggregates.
//
// The score is an arbitary representation of how suitable the Instance is for the
// given set of SortWeights, and is relative to other Instances represented by the
// provided Instance Aggregates.
func CalculateInstanceScoreFromWeights(
	instance *instPkg.Instance,
	aggregates instPkg.Aggregates,
	weights SortWeights,
) float64 {
	normalisedVcpu := aggregates.NormaliseVcpu(instance.Vcpu)
	normalisedRp := aggregates.NormaliseRevocationProbability(instance.RevocationProbability)
	normalisedPrice := aggregates.NormalisePricePerHour(instance.PricePerHour)

	return weights.VcpuWeight*normalisedVcpu +
		weights.RevocationProbabilityWeight*normalisedRp +
		weights.PriceWeight*normalisedPrice
}

// CalculateInstanceScoreFromWeightsWithVcpuLimiter computes a score in the same way as
// CalculateInstanceScoreFromWeights, except a limiter is used to limit the effect of
// the Instance's Vcpu on the score.
func CalculateInstanceScoreFromWeightsWithVcpuLimiter(
	instance *instPkg.Instance,
	aggregates instPkg.Aggregates,
	weights SortWeights,
	maxVcpu int,
) float64 {
	if instance.Vcpu >= maxVcpu {
		return calculatedInstanceScoreFromWeightsWithFixedVcpu(
			instance,
			aggregates,
			weights,
			maxVcpu,
		)
	}
	return CalculateInstanceScoreFromWeights(
		instance,
		aggregates,
		weights,
	)
}

func calculatedInstanceScoreFromWeightsWithFixedVcpu(
	instance *instPkg.Instance,
	aggregates instPkg.Aggregates,
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
