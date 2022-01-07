package schema

// An Advice describes a list of suggested offerings to
// purchase for a given set of services and constraints
// for one or more regions.
type Advice map[string]RegionAdvice

// A RegionAdvice describes a list of suggested offerings to
// purchase for a given set of services and constraints
// for exactly one region.
type RegionAdvice struct {
	Score       float64              `json:"score"`
	Instances   map[string]*Instance `json:"instances"` // ID to instance
	Assignments Assignments          `json:"assignments"`
}

// Assignments lists the relationships between Services and Instances
// within RegionAdvice.
type Assignments struct {
	ServicesToInstances map[string][]string `json:"servicesToInstances"`
	InstancesToServices map[string][]string `json:"instancesToServices"`
}

// GetAssignedInstancesForService returns the list of Instance IDs which are assigned to a Service
// in a RegionAdvice.
func (ra *RegionAdvice) GetAssignedInstancesForService(serviceName string) []*Instance {
	instances := []*Instance{}
	assignedIds := ra.Assignments.ServicesToInstances[serviceName]

	for _, id := range assignedIds {
		inst := ra.Instances[id]
		instances = append(instances, inst)
	}

	return instances
}

// AddAssignment adds the required information to a RegionAdvicce for a Service to be
// considered "assigned" to an Instance and vice versa.
func (ra *RegionAdvice) AddAssignment(serviceName string, instance *Instance) {
	if ra.Instances == nil {
		ra.Instances = make(map[string]*Instance)
	}

	if _, exists := ra.Instances[instance.Id]; !exists {
		ra.Instances[instance.Id] = instance
	}

	ra.Assignments.add(serviceName, instance.Id)
}

func (a *Assignments) add(serviceName string, instanceId string) {
	if a.InstancesToServices == nil {
		a.InstancesToServices = make(map[string][]string)
	}

	a.InstancesToServices[instanceId] = append(
		a.InstancesToServices[instanceId],
		serviceName,
	)

	if a.ServicesToInstances == nil {
		a.ServicesToInstances = make(map[string][]string)
	}

	a.ServicesToInstances[serviceName] = append(
		a.ServicesToInstances[serviceName],
		instanceId,
	)
}
