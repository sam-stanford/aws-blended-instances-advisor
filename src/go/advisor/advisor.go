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

	// TODO: ScoreRegionAdvice or something
}
