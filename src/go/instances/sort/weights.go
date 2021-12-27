package sort

import "aws-blended-instances-advisor/api/schema"

// A SortWeights contains the weights to be used when sorting instances.
type SortWeights struct {
	VcpuWeight                  float64
	RevocationProbabilityWeight float64
	PriceWeight                 float64
}

// NewSortWeightsFromApiWeights creates a SortWeights variable from the api/schema package's
// AdvisorWeights.
func NewSortWeightsFromApiWeights(apiWeights schema.AdvisorWeights) SortWeights {
	return SortWeights{
		VcpuWeight:                  -1.0 * apiWeights.Performance,
		RevocationProbabilityWeight: apiWeights.Availability,
		PriceWeight:                 apiWeights.Price,
	}
}
