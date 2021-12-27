package advisor

import (
	"aws-blended-instances-advisor/api/schema"
	instPkg "aws-blended-instances-advisor/instances"
)

type Advisor interface {
	Advise(
		instancesInfo instPkg.GlobalInfo,
		services []schema.Service,
		options schema.Options,
	) (
		*schema.Advice,
		error,
	)

	AdviseForRegion(
		info instPkg.RegionInfo,
		services []schema.Service,
		options schema.Options,
	) (
		*schema.RegionAdvice,
		error,
	)

	ScoreRegionAdvice(
		advice *schema.RegionAdvice,
		globalAgg instPkg.Aggregates,
		services []schema.Service,
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
