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

// WeightedAdvisor is an Advisor which uses a given set of SortWeights to
// select and score Instances.
type WeightedAdvisor struct {
	weights instSort.SortWeights
}

// NewWeightedAdvisor creates a WeightedAdvisor, converting API schema
// AdvisorWeights into the required format.
func NewWeightedAdvisor(weights schema.AdvisorWeights) Advisor {
	return WeightedAdvisor{
		weights: instSort.NewSortWeightsFromApiWeights(weights),
	}
}

// Advicse selects and scores Instances from a group of available
// Instances for all Regions, returning the selection and information as an
// Advice.
func (advisor WeightedAdvisor) Advise(
	instancesInfo instPkg.GlobalInfo,
	services []schema.Service,
	options schema.Options,
) (
	*schema.Advice,
	error,
) {
	advice := make(schema.Advice)

	awsRegions, err := awsTypes.NewRegions(options.Regions)
	if err != nil {
		return nil, utils.PrependToError(err, "could not parse regions")
	}

	for _, region := range awsRegions {
		info, ok := instancesInfo.RegionInfoMap[region]
		if !ok {
			return nil, fmt.Errorf("region not in map: %s", region.CodeString())
		}

		regionAdvice, err := advisor.AdviseForRegion(info, services, options)
		if err != nil {
			return nil, err
		}

		regionAdvice.Score = advisor.ScoreRegionAdvice(regionAdvice, instancesInfo.GlobalAggregates, services)

		advice[region.CodeString()] = *regionAdvice
	}

	return &advice, nil
}

// AdviseForRegion selects and scores Instances from a group of available
// Instances for one Region, returning the selection and information as a
// RegionAdvice.
func (advisor WeightedAdvisor) AdviseForRegion(
	info instPkg.RegionInfo,
	services []schema.Service,
	options schema.Options,
) (
	*schema.RegionAdvice,
	error,
) {

	permanentInstances := copyInstances(info.PermanentInstances)
	transientInstances := copyInstances(info.TransientInstances)

	if !options.ConsiderFreeInstances {
		permanentInstances = removeFreeInstances(permanentInstances)
		transientInstances = removeFreeInstances(transientInstances)
	}

	allInstances := append(permanentInstances, transientInstances...)

	advice := &schema.RegionAdvice{}

	for _, svc := range services {

		for i := 0; i < svc.MinInstances; i += 1 {
			selectedInstance, err := advisor.selectInstanceForService(
				permanentInstances,
				info.PermanentAggregates,
				svc,
				options,
			)
			if err != nil {
				return nil, err
			}

			advice.AddAssignment(svc.Name, selectedInstance.ToApiSchemaInstance())

			// TODO: Don't seem to be sharing :(
			if options.AvoidRepeatedInstanceTypes {
				permanentInstances = removeInstanceFromSlice(permanentInstances, selectedInstance.Id)
			}
			if options.ShareInstancesBetweenServices {
				permanentInstances = append(permanentInstances, createSharedInstance(selectedInstance, svc))
			}
		}

		remainingInstanceCount := svc.MaxInstances - svc.MinInstances
		for i := 0; i < remainingInstanceCount; i += 1 {
			selectedInstance, err := advisor.selectInstanceForService(
				allInstances,
				info.RegionAggregates,
				svc,
				options,
			)
			if err != nil {
				return nil, err
			}

			advice.AddAssignment(svc.Name, selectedInstance.ToApiSchemaInstance())

			if options.AvoidRepeatedInstanceTypes {
				allInstances = removeInstanceFromSlice(allInstances, selectedInstance.Id)
				if isPermanentInstance(selectedInstance) {
					permanentInstances = removeInstanceFromSlice(permanentInstances, selectedInstance.Id)
				}
			}
			if options.ShareInstancesBetweenServices {
				allInstances = append(allInstances, createSharedInstance(selectedInstance, svc))
				if isPermanentInstance(selectedInstance) {
					permanentInstances = append(permanentInstances, createSharedInstance(selectedInstance, svc))
				}
			}
		}
	}

	return advice, nil
}

func copyInstances(instances []*instPkg.Instance) []*instPkg.Instance {
	new := []*instPkg.Instance{}
	return append(new, instances...)
}

func removeFreeInstances(instances []*instPkg.Instance) []*instPkg.Instance {
	filtered := []*instPkg.Instance{}
	for _, inst := range instances {
		if inst.PricePerHour > 0.0 {
			filtered = append(filtered, inst)
		}
	}
	return filtered
}

func removeInstanceFromSlice(slice []*instPkg.Instance, idToRemove string) []*instPkg.Instance {
	new := []*instPkg.Instance{}
	for _, inst := range slice {
		if inst.Id != idToRemove {
			new = append(new, inst)
		}
	}
	return new
}

func isPermanentInstance(inst *instPkg.Instance) bool {
	return inst.PricePerHour == 0
}

func createSharedInstance(inst *instPkg.Instance, assignedService schema.Service) *instPkg.Instance {
	instanceToShare := inst.MakeCopy()
	instanceToShare.MemoryGb = inst.MemoryGb - assignedService.MinMemory
	instanceToShare.PricePerHour = 0 // Already purchased, so no cost
	return instanceToShare
}

func (advisor WeightedAdvisor) selectInstanceForService(
	instances []*instPkg.Instance,
	aggregates instPkg.Aggregates,
	svc schema.Service,
	options schema.Options,
) (*instPkg.Instance, error) {
	searchStart, searchEnd := 0, len(instances)
	var err error

	// TODO: Different result when using --clear-cache as to when not

	// TODO: Function appears non-deterministic
	// TODO: Function sometimes returns instance with less mem than min mem

	searchStart, err = instSearch.SortAndFindMemory(
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

	return instances[searchStart], nil
}

// ScoreRegionAdvice scores a selection of Instances (as a RegionAdvice),
// returning an arbitrary score.
//
// The returned score can be used to compare RegionAdvices, with higher scores
// meaning a better selection.
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
			revocationProbScore += calculateRevocationProbScore(inst, svc, globalAgg)
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

func calculateRevocationProbScore(inst *schema.Instance, svc schema.Service, agg instPkg.Aggregates) float64 {
	// Min-max scale
	return 1 - ((inst.RevocationProbability - agg.MinRevocationProbability) /
		(agg.MaxRevocationProbability - agg.MinRevocationProbability))
}

func calculatePriceScore(inst *schema.Instance, agg instPkg.Aggregates) float64 {
	// Min-max scale
	return 1 - ((inst.PricePerHour - agg.MinPricePerHour) /
		(agg.MaxPricePerHour - agg.MinPricePerHour))
}
