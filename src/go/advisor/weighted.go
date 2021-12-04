package advisor

import (
	"ec2-test/config"
	instPkg "ec2-test/instances"
	"ec2-test/utils"
	"fmt"
)

// TODO: Docs
// TODO: Make all instance slices pointers

type Weighted struct {
}

func (advisor *Weighted) Advise(
	regionInfoMap instPkg.RegionInfoMap,
	services []config.ServiceDescription,
) (
	[]instPkg.Instance,
	InstanceApplicationMap,
	error,
) {
	for region, info := range regionInfoMap {
		instances, instAppMap, err := advisor.AdviseForRegion(info, services) // TODO var names
		if err != nil {
			return nil, nil, err
		}
		fmt.Printf("\n\nRegion: %s, instances: %+v, instAppMap: %+v\n", region.ToCodeString(), instances, instAppMap) // TODO: Improve
		for _, instance := range instances {
			fmt.Printf("\tInstance: %+v\n\n", *instance)
		}
		// TODO: Calc some form of score
	}
	return nil, nil, nil
}

func (advisor *Weighted) AdviseForRegion(
	info instPkg.RegionInfo,
	services []config.ServiceDescription,
) (
	[]*instPkg.Instance,
	InstanceApplicationMap,
	error,
) {

	selectedInstances := []*instPkg.Instance{}
	instanceApplicationMap := make(InstanceApplicationMap)

	for _, svc := range services {

		// TODO: Need to avoid already used and re-use already suggested instances
		// TODO: Do we need to re-calc aggregates...? Don't think so, but should justify

		// TODO: Func
		for i := 0; i < svc.Instances.MinimumCount; i += 1 {
			selectedInstance, err := advisor.selectInstanceForService(
				info.PermanentInstances.Instances,
				info.PermanentInstances.Aggregates,
				svc,
			)
			if err != nil {
				return nil, nil, err
			}
			selectedInstances = append(selectedInstances, selectedInstance)
			instanceApplicationMap[selectedInstance.Name] = append(
				instanceApplicationMap[selectedInstance.Name],
				svc.Name,
			)
		}

		transientInstanceCount := svc.Instances.TotalCount - svc.Instances.MinimumCount
		for i := 0; i < transientInstanceCount; i += 1 {
			selectedInstance, err := advisor.selectInstanceForService(
				info.AllInstances.Instances,
				info.AllInstances.Aggregates,
				svc,
			)
			if err != nil {
				return nil, nil, err
			}
			selectedInstances = append(selectedInstances, selectedInstance)
			instanceApplicationMap[selectedInstance.Name] = append(
				instanceApplicationMap[selectedInstance.Name],
				svc.Name,
			)
		}

	}

	return selectedInstances, instanceApplicationMap, nil
}

func (advisor *Weighted) selectInstanceForService(
	instances []instPkg.Instance,
	aggregates instPkg.Aggregates,
	svc config.ServiceDescription,
) (*instPkg.Instance, error) {
	searchStart, searchEnd := 0, len(instances)

	// TODO: Different result when using --clear-cache as to when not

	// TODO: Function appears non-determinate
	// TODO: Function sometimes returns instance with less mem than min mem

	searchStart, err := instPkg.SortAndFindMemory(
		instances,
		svc.MinMemory,
		searchStart,
		searchEnd,
	)
	if err != nil {
		return nil, utils.PrependToError(err, "could not find memory in instance slice")
	}

	weights := instPkg.GetSortWeights(svc.GetFocus(), svc.FocusWeight)
	instPkg.SortInstancesWeightedWithVcpuLimiter(
		instances,
		aggregates,
		searchStart,
		searchEnd,
		weights,
		svc.MaxVcpu,
	)

	return &instances[searchStart], nil
}
