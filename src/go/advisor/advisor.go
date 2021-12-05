package advisor

import (
	"ec2-test/api"
	"ec2-test/config"
	"ec2-test/instances"
)

type Advisor interface {
	Advise(
		*instances.RegionInfoMap,
		[]api.Service,
	) (
		api.Advice,
		error,
	)

	AdviseForRegion(
		*instances.RegionInfo,
		*config.Constraints,
	) (
		api.RegionAdvice,
		error,
	)

	// TODO: ScoreRegionAdvice or something
}
