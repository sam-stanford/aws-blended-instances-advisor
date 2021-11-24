package advisor

import (
	awsTypes "ec2-test/aws/types"
	"ec2-test/config"
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
	instances []awsTypes.Instance,
	constraints *config.Constraints,
) (
	[]awsTypes.Instance,
	InstanceApplicationMap, // TODO: Rename to instanceServiceMap
	error,
) {
	for _, service := range constraints.Services {
		// TODO: Abstract out

		// TODO: Sort for each iteration is worse than linear search
		// ... but is necessary for finding desired instances

		// Find min mem
		awsTypes.SortInstancesByMemory(instances, 0, len(instances))
		minMemoryIndex, err := awsTypes.GetIndexOfMinimumMemoryFromSortedInstances(instances, service.MinMemory, 0, len(instances))
		if err != nil {
			return nil, nil, utils.PrependToError(err, "error when finding index of instance with minimum memory requirement")
		}

		// Find non-revocable
		awsTypes.SortInstancesByRevocationProbability(instances, minMemoryIndex, len(instances))
		minRevocationProbabilityIndex, err := awsTypes.GetIndexOfMinimumMemoryFromSortedInstances(instances, 1, minMemoryIndex, len(instances))
		if err != nil {
			return nil, nil, utils.PrependToError(err, "error when finding index of instance with desired revocation probability")
		}

		// Find minimum desired cpu
		awsTypes.SortInstancesByVcpu(instances, minRevocationProbabilityIndex, len(instances))
		desiredVcpuIndex, err := awsTypes.GetIndexOfMinimumVcpuFromSortedInstances(instances, service.MaxVcpu, minRevocationProbabilityIndex, len(instances))
		if err != nil {
			return nil, nil, utils.PrependToError(err, "error when finding index of instance with desired VCPU")
		}

		// Find lowest price
		indexOfLowestPrice := awsTypes.GetIndexOfMinimumPriceFromInstances(instances, -1, desiredVcpuIndex, len(instances)) // -1 for lowest price

	}

	return nil, nil, nil
}

func (advisor NaiveReliabilityAdvisor) AdviseForRegions(
	regionInstancesMap awsTypes.RegionInstancesMap,
	constraints *config.Constraints,
) (
	[]awsTypes.Instance,
	InstanceApplicationMap,
	error,
) {
	for region, instances := range regionInstancesMap {
		advisedInstances, instanceApplicationMap, err := advisor.Advise(instances, constraints)
		fmt.Printf("Region: %s, advisedInstances: %+v, instanceAppMap: %+v, err: %s\n", region.ToCodeString(), advisedInstances, instanceApplicationMap, err.Error())
		// TODO
	}
	// TODO: Score & choose
	return nil, nil, nil
}

// TODO: Maybe a score() function
