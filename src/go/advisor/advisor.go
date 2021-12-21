package advisor

import (
	"aws-blended-instances-advisor/api/schema"
	"aws-blended-instances-advisor/instances"
)

type Advisor interface {
	Advise(
		instances.RegionInfoMap,
		[]schema.Service,
		schema.Options,
	) (
		*schema.Advice,
		error,
	)

	AdviseForRegion(
		instances.RegionInfo,
		[]schema.Service,
	) (
		*schema.RegionAdvice,
		error,
	)

	ScoreRegionAdvice(
		*schema.RegionAdvice,
		instances.Aggregates,
		[]schema.Service,
	) float64
}

// TODO: Take logger in args & log stuff

// New creates an advisor, using the type provided in the info argument
// to determine which advisor to use.
func New(info schema.Advisor) Advisor {
	switch info.Type {
	case schema.Weighted:
		return NewWeightedAdvisor(info.Weights)

	case schema.Random:
		return NewRandomAdvisor()

	default:
		return NewWeightedAdvisor(info.Weights)
	}
}
