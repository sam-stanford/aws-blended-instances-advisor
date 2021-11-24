package advisor

import (
	awsTypes "ec2-test/aws/types"
	"ec2-test/config"
)

type Advisor interface {
	AdviseForRegion(awsTypes.RegionInstancesMap, *config.Constraints) ([]awsTypes.Instance, InstanceApplicationMap, error)
	Advise([]awsTypes.Instance, *config.Constraints) ([]awsTypes.Instance, InstanceApplicationMap, error)
}
