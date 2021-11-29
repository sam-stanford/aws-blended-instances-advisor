package advisor

import (
	"ec2-test/config"
	instances "ec2-test/instance"
)

type Advisor interface {
	Advise(
		[]instances.Instance,
		*config.Constraints,
	) (
		[]instances.Instance,
		InstanceApplicationMap,
		error,
	)

	AdviseForEachRegion(
		instances.RegionInstancesMap,
		*config.Constraints,
	) (
		[]instances.Instance,
		InstanceApplicationMap,
		error,
	)
}
