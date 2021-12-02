package instances

import (
	"ec2-test/config"
	"sort"
	"strings"
)

// TODO: Own sub package?

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
	aggregates Aggregates,
	startIndex,
	endIndex int,
	weightings SortWeightings,
) {
	sort.Slice(instances[startIndex:endIndex], func(i, j int) bool {
		iScore := CalculateInstanceScoreFromWeight(
			instances[startIndex+i],
			aggregates,
			weightings,
		)
		jScore := CalculateInstanceScoreFromWeight(
			instances[startIndex+j],
			aggregates,
			weightings,
		)
		return iScore < jScore
	})

}

// TODO: Doc & test
func CalculateInstanceScoreFromWeight(
	instance Instance,
	aggregates Aggregates,
	weightings SortWeightings,
) float64 {
	normalisedVcpu := aggregates.NormaliseVcpu(instance.Vcpu)
	normalisedRp := aggregates.NormaliseRevocationProbability(instance.RevocationProbability)
	normalisedPrice := aggregates.NormalisePricePerHour(instance.PricePerHour)

	return weightings.VcpuWeight*normalisedVcpu +
		weightings.RevocationProbabilityWeight*normalisedRp +
		weightings.PriceWeight*normalisedPrice
}

// TODO: Doc & test
func GetSortWeights(focus config.ServiceFocus, focusWeight float64) SortWeightings {
	primaryFocusWeight := 0.33 + 2.0*0.33*focusWeight
	secondaryFocusWeight := 0.33 * (1.0 - focusWeight)

	// TODO: Comment on negative weight for vcpu (want to max, while others are min)

	switch focus {
	case config.Availability:
		return SortWeightings{
			RevocationProbabilityWeight: primaryFocusWeight,
			VcpuWeight:                  -1.0 * secondaryFocusWeight,
			PriceWeight:                 secondaryFocusWeight,
		}
	case config.Cost:
		return SortWeightings{
			RevocationProbabilityWeight: secondaryFocusWeight,
			VcpuWeight:                  -1.0 * secondaryFocusWeight,
			PriceWeight:                 primaryFocusWeight,
		}

	case config.Performance:
		return SortWeightings{
			RevocationProbabilityWeight: secondaryFocusWeight,
			VcpuWeight:                  -1.0 * primaryFocusWeight,
			PriceWeight:                 secondaryFocusWeight,
		}
	default: // Implicitly includes config.Balanced
		return SortWeightings{
			RevocationProbabilityWeight: 0.33,
			VcpuWeight:                  0.33,
			PriceWeight:                 0.33,
		}
	}

}
