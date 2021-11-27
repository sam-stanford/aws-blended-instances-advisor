package advisor

import (
	"ec2-test/config"
	instancesPkg "ec2-test/instances"
	"ec2-test/utils"
	"fmt"
)

//
type NaiveReliabilityAdvisor struct {
}

// Instantiates a NaiveReliabilityAdvisor
func NewNaiveReliabilityAdvisor() NaiveReliabilityAdvisor {
	// TODO
	return NaiveReliabilityAdvisor{}
}

// TODO: Doc
func (advisor NaiveReliabilityAdvisor) Advise(
	instances []instancesPkg.Instance,
	constraints *config.Constraints,
) (
	[]instancesPkg.Instance,
	InstanceApplicationMap, // TODO: Rename to instanceServiceMap
	error,
) {

	selectedInstances := []instancesPkg.Instance{}

	for _, service := range constraints.Services {
		// TODO: Abstract out sort & find

		searchStart, searchEnd := 0, len(instances)

		// Find min mem
		instancesPkg.SortInstancesByMemory(instances, 0, len(instances))
		searchStart, err := instancesPkg.FindMinimumMemorySortedInstances(instances, service.MinMemory, searchStart, searchEnd)
		if err != nil {
			return nil, nil, utils.PrependToError(err, "error when finding index of instance with minimum memory requirement")
		}

		// Find non-revocable
		instancesPkg.SortInstancesByRevocationProbability(instances, searchStart, searchEnd)
		searchEnd, err = instancesPkg.FindMinimumRevocationProbabilitySortedInstances(instances, 1, searchStart, searchEnd) // TODO: Change 1 to 0
		if err != nil {
			return nil, nil, utils.PrependToError(err, "error when finding index of instance with desired revocation probability")
		}

		// Find minimum desired cpu
		instancesPkg.SortInstancesByVcpus(instances, searchStart, searchEnd)
		searchStart, err = instancesPkg.FindMinimumVcpuSortedInstances(instances, service.MaxVcpu, searchStart, searchEnd)
		if err != nil {
			return nil, nil, utils.PrependToError(err, "error when finding index of instance with desired VCPU")
		}

		// Find lowest price
		instancesPkg.SortInstancesByPrice(instances, searchStart, searchEnd)
		selectedInstance := instances[searchStart]

		selectedInstances = append(selectedInstances, selectedInstance)
	}

	return selectedInstances, nil, nil
}

func (advisor NaiveReliabilityAdvisor) AdviseForRegions(
	regionInstancesMap instancesPkg.RegionInstancesMap,
	constraints *config.Constraints,
) (
	[]instancesPkg.Instance,
	InstanceApplicationMap,
	error,
) {
	for region, instances := range regionInstancesMap {
		advisedInstances, instanceApplicationMap, err := advisor.Advise(instances, constraints)
		fmt.Println(region, advisedInstances, instanceApplicationMap, err)
		// TODO
	}
	// TODO: Score & choose
	return nil, nil, nil
}

// TODO: Maybe a score() function
