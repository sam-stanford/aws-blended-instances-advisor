package advisor

import (
	"ec2-test/config"
	"ec2-test/instance"
	"ec2-test/utils"
)

type LowestCostAdvisor struct {
}

func (advisor *LowestCostAdvisor) Advise(
	instances []instance.Instance,
	services []config.ServiceDescription,
) (
	[]instance.Instance,
	InstanceApplicationMap,
	error,
) {

	selectedInstances := []instance.Instance{}
	instanceApplicationMap := make(InstanceApplicationMap)

	searchStart, searchEnd := 0, len(instances)

	for _, svc := range services {
		searchStart, err := instance.SortAndFindMemory(instances, svc.MinMemory, searchStart, searchEnd)
		if err != nil {
			return nil, nil, utils.PrependToError(err, "could not find memory in instance slice")
		}

		instance.SortInstancesByPrice(instances, searchStart, searchEnd)

		// TODO: Abstract out
		currentIndex := searchStart
		lowestPrice := instances[searchStart].PricePerHour
		for currentIndex < searchEnd && instances[currentIndex].PricePerHour == lowestPrice {
			currentIndex += 1
		}
		searchEnd = currentIndex - 1

		// TODO: Here, but maybe need to do the weighted ranking thing instead

	}
}

func (advisor *LowestCostAdvisor) AdvisePermanentInstances(instances []instance.Instance)

func (advisor *LowestCostAdvisor) AdviseTransientInstances(instances []instance.Instance)

func (advisor *LowestCostAdvisor) AdviseForEachRegion(
	regionInstanceMap instance.RegionInstancesMap,
	services []config.ServiceDescription,
) (
	[]instance.Instance,
	InstanceApplicationMap,
	error,
) {
	// TODO
	return nil, nil, nil
}
