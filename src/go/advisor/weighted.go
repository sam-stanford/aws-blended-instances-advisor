package advisor

import (
	"ec2-test/api"
	instPkg "ec2-test/instances"
	"ec2-test/utils"
)

// TODO: Docs
// TODO: Make all instance slices pointers
// TODO: Rename to FocusAdvisor & create CustomAdvisor which takes specified weights

type Weighted struct {
	focus       api.AdvisorFocus
	focusWeight float64
}

func NewWeighted(focus api.AdvisorFocus, focusWeight float64) Weighted {
	return Weighted{
		focus:       focus,
		focusWeight: focusWeight,
	}
}

func (advisor Weighted) Advise(
	regionInfoMap instPkg.RegionInfoMap,
	services []api.Service,
) (
	*api.Advice,
	error,
) {
	advice := make(api.Advice)

	for region, info := range regionInfoMap {
		regionAdvice, err := advisor.AdviseForRegion(info, services)
		if err != nil {
			return nil, err
		}

		advice[region.ToCodeString()] = *regionAdvice
	}

	return &advice, nil
}

func (advisor Weighted) AdviseForRegion(
	info instPkg.RegionInfo,
	services []api.Service,
) (
	*api.RegionAdvice,
	error,
) {
	advice := &api.RegionAdvice{}

	for _, svc := range services {

		// TODO: Need to avoid already used and re-use already suggested instances
		// TODO: Do we need to re-calc aggregates...? Don't think so, but should justify

		// TODO: Func (repeated code)
		for i := 0; i < svc.MinInstances; i += 1 {
			selectedInstance, err := advisor.selectInstanceForService(
				info.PermanentInstances.Instances,
				info.PermanentInstances.Aggregates,
				svc,
			)
			if err != nil {
				return nil, err
			}
			advice.AddAssignment(svc.Name, selectedInstance)
		}

		transientInstances := svc.TotalInstances - svc.MinInstances
		for i := 0; i < transientInstances; i += 1 {
			selectedInstance, err := advisor.selectInstanceForService(
				info.AllInstances.Instances,
				info.AllInstances.Aggregates,
				svc,
			)
			if err != nil {
				return nil, err
			}
			advice.AddAssignment(svc.Name, selectedInstance)
		}
	}

	// TODO: Calc some form of score

	return advice, nil
}

func (advisor Weighted) selectInstanceForService(
	instances []*instPkg.Instance,
	aggregates instPkg.Aggregates,
	svc api.Service,
) (*api.Instance, error) {
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

	// TODO: Allow for user to pick these - maybe a diff advisor
	weights := instPkg.GetSortWeights(advisor.focus, advisor.focusWeight)
	instPkg.SortInstancesWeightedWithVcpuLimiter(
		instances,
		aggregates,
		searchStart,
		searchEnd,
		weights,
		svc.MaxVcpu,
	)

	return instances[searchStart].ToApiInstance(), nil
}
