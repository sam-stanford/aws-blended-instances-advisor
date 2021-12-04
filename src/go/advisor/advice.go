package advisor

import "ec2-test/instances"

type Advice struct {
	Instances   []instances.Instance
	Assignments InstanceApplicationMap
}

// TODO: Use Name+Region or Name+AZ for key
// TODO: Get & set methods
type InstanceApplicationMap map[string][]string
