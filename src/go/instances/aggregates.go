package instances

import (
	"ec2-test/utils"
)

// Aggregates contains aggregate information for a given set of instances,
// providing the ability to normalise/standardise properties of an instance.
type Aggregates struct {
	Count int

	MinVcpu  int
	MaxVcpu  int
	MeanVcpu float64

	MinRevocationProbability  float64
	MaxRevocationProbability  float64
	MeanRevocationProbability float64

	MinPricePerHour  float64
	MaxPricePerHour  float64
	MeanPricePerHour float64
}

// Calculates aggregates for a slice of instances, returning information in
// an Aggregate struct.
func CalculateAggregates(instances []*Instance) Aggregates {
	totalVcpu, minVcpu, maxVcpu := 0, instances[0].Vcpu, instances[0].Vcpu

	totalPricePerHour := 0.0
	minPricePerHour, maxPricePerHour := instances[0].PricePerHour, instances[0].PricePerHour

	totalRevocationProbability := 0.0
	minRevocationProbability := instances[0].RevocationProbability
	maxRevocationProbability := instances[0].RevocationProbability

	for _, instance := range instances {
		totalVcpu += instance.Vcpu
		totalRevocationProbability += instance.RevocationProbability
		totalPricePerHour += instance.PricePerHour

		minVcpu = utils.MinOfInts(minVcpu, instance.Vcpu)
		minRevocationProbability = utils.MinOfFloats(minRevocationProbability, instance.RevocationProbability)
		minPricePerHour = utils.MinOfFloats(minPricePerHour, instance.PricePerHour)

		maxVcpu = utils.MaxOfInts(maxVcpu, instance.Vcpu)
		maxRevocationProbability = utils.MaxOfFloats(maxRevocationProbability, instance.RevocationProbability)
		maxPricePerHour = utils.MaxOfFloats(maxPricePerHour, instance.PricePerHour)
	}

	floatCount := float64(len(instances))
	return Aggregates{
		Count: len(instances),

		MeanVcpu:                  float64(totalVcpu) / floatCount,
		MeanRevocationProbability: totalRevocationProbability / floatCount,
		MeanPricePerHour:          totalPricePerHour / floatCount,

		MinVcpu:                  minVcpu,
		MinRevocationProbability: minRevocationProbability,
		MinPricePerHour:          minPricePerHour,

		MaxVcpu:                  maxVcpu,
		MaxRevocationProbability: maxRevocationProbability,
		MaxPricePerHour:          maxPricePerHour,
	}
}

// Normalises a given VCPU value with respect to aggregate values using min-max scaling.
// Returns 1/count if aggregates are formed from all equal values.
func (agg Aggregates) NormaliseVcpu(vcpu int) float64 {
	if agg.MaxVcpu == agg.MinVcpu {
		return 1.0 / float64(agg.Count)
	}
	return float64(vcpu-agg.MinVcpu) / float64(agg.MaxVcpu-agg.MinVcpu)
}

// Normalises a given RevocationProbablity with respect to aggregate values using min-max scaling.
// Returns 1/count if aggregates are formed from all equal values.
func (agg Aggregates) NormaliseRevocationProbability(prob float64) float64 {
	if utils.FloatsEqual(agg.MaxRevocationProbability, agg.MinRevocationProbability) {
		return 1.0 / float64(agg.Count)
	}
	return (prob - agg.MinRevocationProbability) / (agg.MaxRevocationProbability - agg.MinRevocationProbability)
}

// Normalises a given PricePerHour value with respect to aggregate values using min-max scaling.
// Returns 1/count if aggregates are formed from all equal values.
func (agg Aggregates) NormalisePricePerHour(price float64) float64 {
	if utils.FloatsEqual(agg.MaxPricePerHour, agg.MinPricePerHour) {
		return 1.0 / float64(agg.Count)
	}
	return (price - agg.MinPricePerHour) / (agg.MaxPricePerHour - agg.MinPricePerHour)
}
