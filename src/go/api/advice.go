package api

// TODO: Doc in README & point to README from here

type Advice map[string]RegionAdvice

type RegionAdvice struct {
	Score       float64              `json:"score"`
	Instances   map[string]*Instance `json:"instances"` // ID to instance
	Assignments Assignments          `json:"assignments"`
}

type Assignments struct {
	ServicesToInstances map[string][]string `json:"servicesToInstances"`
	InstancesToServices map[string][]string `json:"instancesToServices"`
}

// TODO: Doc & test
func (ra *RegionAdvice) GetAssignedInstancesForService(serviceName string) []*Instance {
	instances := []*Instance{}
	assignedIds := ra.Assignments.ServicesToInstances[serviceName]

	for _, id := range assignedIds {
		inst := ra.Instances[id]
		instances = append(instances, inst)
	}

	return instances
}

// TODO: Use "NewAdvice" and "NewRegionAdvice" to instantiate maps rather than checking on each access

// TODO: Doc & test
func (ra *RegionAdvice) AddAssignment(serviceName string, instance *Instance) {
	if ra.Instances == nil {
		ra.Instances = map[string]*Instance{
			instance.Id: instance,
		}
	} else {
		ra.Instances[instance.Id] = instance
	}

	ra.Assignments.add(serviceName, instance.Id)
}

func (a *Assignments) add(serviceName string, instanceId string) {
	if a.InstancesToServices == nil {
		a.InstancesToServices = map[string][]string{
			instanceId: {serviceName},
		}
	} else {
		a.InstancesToServices[instanceId] = appendStringIfNotInSlice(a.InstancesToServices[instanceId], serviceName)
	}

	if a.ServicesToInstances == nil {
		a.ServicesToInstances = map[string][]string{
			serviceName: {instanceId},
		}
	} else {
		a.ServicesToInstances[serviceName] = appendStringIfNotInSlice(a.ServicesToInstances[serviceName], instanceId)
	}
}

// TODO: Move to utils & test
func appendStringIfNotInSlice(slice []string, s string) []string {
	inSlice := false
	for i := range slice {
		if slice[i] == s {
			inSlice = true
			break
		}
	}

	if !inSlice {
		return append(slice, s)
	}
	return slice
}
