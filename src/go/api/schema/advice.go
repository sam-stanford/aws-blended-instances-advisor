package schema

import "aws-blended-instances-advisor/utils"

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
		a.InstancesToServices[instanceId] = utils.AppendStringIfNotInSlice(a.InstancesToServices[instanceId], serviceName)
	}

	if a.ServicesToInstances == nil {
		a.ServicesToInstances = map[string][]string{
			serviceName: {instanceId},
		}
	} else {
		a.ServicesToInstances[serviceName] = utils.AppendStringIfNotInSlice(a.ServicesToInstances[serviceName], instanceId)
	}
}
