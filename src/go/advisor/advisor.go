package advisor

import (
	"ec2-test/api"
	"ec2-test/instances"
)

type Advisor interface {
	Advise(
		instances.RegionInfoMap,
		[]api.Service,
	) (
		*api.Advice,
		error,
	)

	AdviseForRegion(
		instances.RegionInfo,
		[]api.Service,
	) (
		*api.RegionAdvice,
		error,
	)

	ScoreRegionAdvice(
		*api.RegionAdvice,
		instances.Aggregates,
		[]api.Service,
	) float64
}

// TODO: Take logger in args & log stuff

// New creates an advisor, using the type provided in the info argument
// to determine which advisor to use.
func New(info api.Advisor) Advisor {
	switch info.Type {
	case api.Weighted:
		return NewWeightedAdvisor(info.Focus, info.FocusWeight)

		// TODO: Random, custom, etc.

	default:
		// TODO: Which default? Maybe naive or something
		return NewWeightedAdvisor(info.Focus, info.FocusWeight)
	}
}
