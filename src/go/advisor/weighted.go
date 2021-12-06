package advisor

import (
	"ec2-test/api"
	instPkg "ec2-test/instances"
	"ec2-test/utils"
)

// TODO: Docs
// TODO: Make all instance slices pointers
// TODO: Rename to FocusAdvisor & create CustomAdvisor which takes specified weights
// TODO: Use logger

type WeightedAdvisor struct {
	focus       api.AdvisorFocus
	focusWeight float64
	weights     instPkg.SortWeightings
}

func NewWeightedAdvisor(focus api.AdvisorFocus, focusWeight float64) WeightedAdvisor {
	return WeightedAdvisor{
		focus:       focus,
		focusWeight: focusWeight,
		weights:     instPkg.GetSortWeights(focus, focusWeight), // TODO: Rename this to be FocusAdvisor & have WeightedAdvisor be Custom
	}
}

func (advisor WeightedAdvisor) Advise(
	regionInfoMap instPkg.RegionInfoMap,
	services []api.Service,
) (
	*api.Advice,
	error,
) {
	advice := make(api.Advice)
	globalAggregates := getGlobalAggregatesFromRegionInfoMap(regionInfoMap)

	for region, info := range regionInfoMap {
		regionAdvice, err := advisor.AdviseForRegion(info, services)
		if err != nil {
			return nil, err
		}

		regionAdvice.Score = advisor.ScoreRegionAdvice(regionAdvice, globalAggregates, services)

		advice[region.ToCodeString()] = *regionAdvice
	}

	return &advice, nil
}

// TODO: Move to inst package?
func getGlobalAggregatesFromRegionInfoMap(m instPkg.RegionInfoMap) instPkg.Aggregates {
	aggs := []instPkg.Aggregates{}
	for _, info := range m {
		aggs = append(aggs, info.AllInstances.Aggregates)
	}
	return instPkg.CombineAggregates(aggs)
}

func (advisor WeightedAdvisor) AdviseForRegion(
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

func (advisor WeightedAdvisor) selectInstanceForService(
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

	instPkg.SortInstancesWeightedWithVcpuLimiter(
		instances,
		aggregates,
		searchStart,
		searchEnd,
		advisor.weights,
		svc.MaxVcpu,
	)

	return instances[searchStart].ToApiInstance(), nil
}

// TODO: Test & Doc
func (advisor WeightedAdvisor) ScoreRegionAdvice(
	advice *api.RegionAdvice,
	globalAgg instPkg.Aggregates,
	services []api.Service,
) float64 {
	vcpuScore, revocationProbScore, priceScore := 0.0, 0.0, 0.0
	totalInstances := 0

	for _, svc := range services {
		assignedInstances := advice.GetAssignedInstancesForService(svc.Name)
		for _, inst := range assignedInstances {
			vcpuScore += calculateVcpuScore(inst, svc)
			revocationProbScore += calculateRevocationProbScore(inst, globalAgg)
			priceScore += calculatePriceScore(inst, globalAgg)
		}
		totalInstances += len(assignedInstances)
	}

	return ((vcpuScore * advisor.weights.VcpuWeight) +
		(revocationProbScore * advisor.weights.RevocationProbabilityWeight) +
		(priceScore * advisor.weights.PriceWeight)) / float64(totalInstances)
}

func calculateVcpuScore(inst *api.Instance, svc api.Service) float64 {
	// Percentage of MaxVcpu
	if inst.Vcpu >= svc.MaxVcpu {
		return 1.0
	}
	return float64(inst.Vcpu) / float64(svc.MaxVcpu)
}

func calculateRevocationProbScore(inst *api.Instance, agg instPkg.Aggregates) float64 {
	// Min-max scale
	return 1 - ((inst.RevocationProbability - agg.MinRevocationProbability) /
		(agg.MaxRevocationProbability - agg.MinRevocationProbability))
}

func calculatePriceScore(inst *api.Instance, agg instPkg.Aggregates) float64 {
	// Min-max scale
	return 1 - ((inst.PricePerHour - agg.MinPricePerHour) /
		(agg.MaxPricePerHour - agg.MinPricePerHour))
}
