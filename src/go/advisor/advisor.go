package advisor

import (
	"ec2-test/aws/types"
	"ec2-test/config"
	"ec2-test/instances"
)

type Advisor interface {
	AdviseForRegion(
		*instances.RegionInfo,
		*config.Constraints,
	) (
		Advice,
		error,
	)

	Advise(
		*instances.RegionInfoMap,
		*config.Constraints,
	) (
		map[types.Region]Advice,
		error,
	)

	// TODO: Make Advise generic and create a AdviseForOneRegion
	// TODO: ... and have Advise return one Advice selection
}
