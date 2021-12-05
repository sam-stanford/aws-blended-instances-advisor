package api

// TODO: Doc in README & point to README from here

type Advice map[string]RegionAdvice

type RegionAdvice struct {
	Score       float64     `json:"score"`
	Instances   []Instance  `json:"instances"`
	Assignments Assignments `json:"assignments"`
}

type Assignments struct {
	ServicesToInstances map[string][]string `json:"servicesToInstances"`
	InstancesToServices map[string][]string `json:"instancesToServices"`
}

// TODO: Doc & test
func (ra *RegionAdvice) AddAssignment(serviceName string, instance *Instance) {
	ra.Instances = appendInstanceIfNotInSlice(ra.Instances, instance)
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

func appendInstanceIfNotInSlice(slice []Instance, instance *Instance) []Instance {
	inSlice := false
	for i := range slice {
		if slice[i].Id == instance.Id {
			inSlice = true
			break
		}
	}

	if !inSlice {
		return append(slice, *instance)
	}
	return slice
}

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
