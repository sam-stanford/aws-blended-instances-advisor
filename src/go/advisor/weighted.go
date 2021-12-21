package advisor

import (
	"aws-blended-instances-advisor/api/schema"
	awsTypes "aws-blended-instances-advisor/aws/types"
	instPkg "aws-blended-instances-advisor/instances"
	instSearch "aws-blended-instances-advisor/instances/search"
	instSort "aws-blended-instances-advisor/instances/sort"
	"aws-blended-instances-advisor/utils"
	"fmt"
)

// TODO: Docs
// TODO: Use logger

type WeightedAdvisor struct {
	weights instSort.SortWeights
}

func NewWeightedAdvisor(weights schema.AdvisorWeights) Advisor {
	return WeightedAdvisor{
		weights: instSort.NewSortWeightsFromApiWeights(weights),
	}
}

func (advisor WeightedAdvisor) Advise(
	regionInfoMap instPkg.RegionInfoMap,
	services []schema.Service,
	options schema.Options,
) (
	*schema.Advice,
	error,
) {
	advice := make(schema.Advice) // TODO: Use NewAdvice here
	globalAggregates := getGlobalAggregatesFromRegionInfoMap(regionInfoMap)

	awsRegions, err := awsTypes.NewRegions(options.Regions)
	if err != nil {
		return nil, utils.PrependToError(err, "could not parse regions")
	}

	for _, region := range awsRegions {
		info, ok := regionInfoMap[region]
		if !ok {
			return nil, fmt.Errorf("region not in map: %s", region.CodeString())
		}

		regionAdvice, err := advisor.AdviseForRegion(info, services)
		if err != nil {
			return nil, err
		}

		regionAdvice.Score = advisor.ScoreRegionAdvice(regionAdvice, globalAggregates, services)

		advice[region.CodeString()] = *regionAdvice
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
	services []schema.Service,
) (
	*schema.RegionAdvice,
	error,
) {
	advice := &schema.RegionAdvice{}

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
	svc schema.Service,
) (*schema.Instance, error) {
	searchStart, searchEnd := 0, len(instances)

	// TODO: Different result when using --clear-cache as to when not

	// TODO: Function appears non-deterministic
	// TODO: Function sometimes returns instance with less mem than min mem

	searchStart, err := instSearch.SortAndFindMemory(
		instances,
		svc.MinMemory,
		searchStart,
		searchEnd,
	)
	if err != nil {
		return nil, utils.PrependToError(err, "could not find memory in instance slice")
	}

	instSort.SortInstancesWeightedWithVcpuLimiter(
		instances,
		aggregates,
		searchStart,
		searchEnd,
		advisor.weights,
		svc.MaxVcpu,
	)

	return instances[searchStart].ToApiSchemaInstance(), nil
}

// TODO: Test & Doc
func (advisor WeightedAdvisor) ScoreRegionAdvice(
	advice *schema.RegionAdvice,
	globalAgg instPkg.Aggregates,
	services []schema.Service,
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

func calculateVcpuScore(inst *schema.Instance, svc schema.Service) float64 {
	// Percentage of MaxVcpu
	if inst.Vcpu >= svc.MaxVcpu {
		return 1.0
	}
	return float64(inst.Vcpu) / float64(svc.MaxVcpu)
}

func calculateRevocationProbScore(inst *schema.Instance, agg instPkg.Aggregates) float64 {
	// Min-max scale
	return 1 - ((inst.RevocationProbability - agg.MinRevocationProbability) /
		(agg.MaxRevocationProbability - agg.MinRevocationProbability))
}

func calculatePriceScore(inst *schema.Instance, agg instPkg.Aggregates) float64 {
	// Min-max scale
	return 1 - ((inst.PricePerHour - agg.MinPricePerHour) /
		(agg.MaxPricePerHour - agg.MinPricePerHour))
}
