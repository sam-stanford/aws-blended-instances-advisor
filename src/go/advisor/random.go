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

func (a RandomAdvisor) ScoreRegionAdvice(
	advice *schema.RegionAdvice,
	globalAgg instPkg.Aggregates,
	services []schema.Service,
) float64 {
	// TODO
	return 0
}
