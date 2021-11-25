package advisor

import (
	"ec2-test/config"
	"ec2-test/instances"
)

type Advisor interface {
	Advise([]instances.Instance, *config.Constraints) ([]instances.Instance, InstanceApplicationMap, error)
	AdviseForRegions(instances.RegionInstancesMap, *config.Constraints) ([]instances.Instance, InstanceApplicationMap, error)
}
