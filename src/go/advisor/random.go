package advisor

import (
	"aws-blended-instances-advisor/api/schema"
	instPkg "aws-blended-instances-advisor/instances"
)

type RandomAdvisor struct{}

func NewRandomAdvisor() Advisor {
	return RandomAdvisor{}
}

func (a RandomAdvisor) Advise(
	instances instPkg.RegionInfoMap,
	services []schema.Service,
	options schema.Options,
) (
	*schema.Advice,
	error,
) {
	// TODO
	return nil, nil
}

func (a RandomAdvisor) AdviseForRegion(
	regionInfo instPkg.RegionInfo,
	services []schema.Service,
) (
	*schema.RegionAdvice,
	error,
) {
	// TODO
	return nil, nil
}

func (a RandomAdvisor) ScoreRegionAdvice(
	advice *schema.RegionAdvice,
	globalAggregates instPkg.Aggregates,
	services []schema.Service,
) float64 {
	// TODO
	return 0
}
