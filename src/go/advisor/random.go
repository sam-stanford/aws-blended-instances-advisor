package advisor

import (
	"aws-blended-instances-advisor/api/schema"
	instPkg "aws-blended-instances-advisor/instances"
)

type RandomAdvisor struct{}

func NewRandomAdvisor() Advisor {
	return RandomAdvisor{}
}

// Advicse selects and scores Instances from a group of available
// Instances for all Regions, returning the selection and information as an
// Advice.
func (a RandomAdvisor) Advise(
	instancesInfo instPkg.GlobalInfo,
	services []schema.Service,
	options schema.Options,
) (
	*schema.Advice,
	error,
) {
	// TODO
	return nil, nil
}

// AdviseForRegion selects and scores Instances from a group of available
// Instances for one Region, returning the selection and information as a
// RegionAdvice.
func (a RandomAdvisor) AdviseForRegion(
	info instPkg.RegionInfo,
	services []schema.Service,
	options schema.Options,
) (
	*schema.RegionAdvice,
	error,
) {
	// TODO
	return nil, nil
}

// ScoreRegionAdvice scores a selection of Instances (as a RegionAdvice),
// returning an arbitrary score.
//
// The returned score can be used to compare RegionAdvices, with higher scores
// meaning a better selection.
func (a RandomAdvisor) ScoreRegionAdvice(
	advice *schema.RegionAdvice,
	globalAgg instPkg.Aggregates,
	services []schema.Service,
) float64 {
	// TODO
	return 0
}
