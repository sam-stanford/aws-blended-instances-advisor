package advisor

// TODO
type Naive struct {
}

// TODO: Maybe "random" instead?

// TODO: Simply choose the highest ranking weight & use that

// import (
// 	"ec2-test/config"
// 	"ec2-test/instance"
// 	"ec2-test/utils"
// 	"fmt"
// )

// // TODO: Doc
// type NaivePriceAdvisor struct {
// }

// // Instantiates a NaivePriceAdvisor
// func NewNaivePriceAdvisor() NaiveReliabilityAdvisor {
// 	return NaiveReliabilityAdvisor{}
// }

// // TODO: Doc
// func (advisor NaiveReliabilityAdvisor) Advise(
// 	instances []instance.Instance,
// 	constraints *config.Constraints,
// ) (
// 	[]instance.Instance,
// 	InstanceApplicationMap, // TODO: Rename to instanceServiceMap
// 	error,
// ) {

// 	selectedInstances := []instance.Instance{}

// 	for _, service := range constraints.Services {
// 		// TODO: Abstract out sort & find

// 		searchStart, searchEnd := 0, len(instances)

// 		// Find min mem
// 		instance.SortInstancesByMemory(instances, 0, len(instances))
// 		searchStart, err := instance.FindMinimumMemorySortedInstances(instances, service.MinMemory, searchStart, searchEnd)
// 		if err != nil {
// 			return nil, nil, utils.PrependToError(err, "error when finding index of instance with minimum memory requirement")
// 		}

// 		// Find non-revocable
// 		instance.SortInstancesByRevocationProbability(instances, searchStart, searchEnd)
// 		searchEnd, err = instance.FindMinimumRevocationProbabilitySortedInstances(instances, 1, searchStart, searchEnd) // TODO: Change 1 to 0
// 		if err != nil {
// 			return nil, nil, utils.PrependToError(err, "error when finding index of instance with desired revocation probability")
// 		}

// 		// Find minimum desired cpu
// 		instance.SortInstancesByVcpu(instances, searchStart, searchEnd)
// 		searchStart, err = instance.FindMinimumVcpuSortedInstances(instances, service.MaxVcpu, searchStart, searchEnd)
// 		if err != nil {
// 			return nil, nil, utils.PrependToError(err, "error when finding index of instance with desired VCPU")
// 		}

// 		// Find lowest price
// 		instance.SortInstancesByPrice(instances, searchStart, searchEnd)
// 		selectedInstance := instances[searchStart]

// 		selectedInstances = append(selectedInstances, selectedInstance)
// 	}

// 	return selectedInstances, nil, nil
// }

// func (advisor NaiveReliabilityAdvisor) AdviseForEachRegion(
// 	regionInstancesMap instance.RegionInstancesMap,
// 	constraints *config.Constraints,
// ) (
// 	[]instance.Instance,
// 	InstanceApplicationMap,
// 	error,
// ) {
// 	for region, instances := range regionInstancesMap {
// 		advisedInstances, instanceApplicationMap, err := advisor.Advise(instances, constraints)
// 		fmt.Println(region, advisedInstances, instanceApplicationMap, err)
// 		// TODO
// 	}
// 	// TODO: Score & choose
// 	return nil, nil, nil
// }

// // TODO: Maybe a score() function
