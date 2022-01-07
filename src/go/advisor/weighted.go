package advisor

import (
	"aws-blended-instances-advisor/api/schema"
	awsTypes "aws-blended-instances-advisor/aws/types"
	instPkg "aws-blended-instances-advisor/instances"
	instSearch "aws-blended-instances-advisor/instances/search"
	instSort "aws-blended-instances-advisor/instances/sort"
	"aws-blended-instances-advisor/utils"
	"fmt"

	"go.uber.org/zap"
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
	logger *zap.Logger,
) (
	*schema.Advice,
	error,
) {
	logger.Info(
		"advising with weighted advisor",
		zap.Any("weights", advisor.weights),
	)

	advice := make(schema.Advice)

	awsRegions, err := awsTypes.NewRegions(options.Regions)
	if err != nil {
		return nil, utils.PrependToError(err, "could not parse regions")
	}

	for _, region := range awsRegions {
		logger.Info("advising for region", zap.String("region", region.CodeString()))

		info, ok := instancesInfo.RegionInfoMap[region]
		if !ok {
			return nil, fmt.Errorf("region not in map: %s", region.CodeString())
		}

		regionAdvice, err := advisor.AdviseForRegion(info, services, options, logger)
		if err != nil {
			return nil, err
		}

		regionAdvice.Score = advisor.ScoreRegionAdvice(regionAdvice, instancesInfo.GlobalAggregates, services, logger)

		advice[region.CodeString()] = *regionAdvice

		logger.Info(
			"advice created for region",
			zap.String("region", region.CodeString()),
			zap.Any("advice", regionAdvice),
		)
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
	logger *zap.Logger,
) (
	*schema.RegionAdvice,
	error,
) {
	logger.Info(
		"advising for region with weighted advisor",
		zap.Any("weights", advisor.weights),
	)

	permanentInstances := copyInstances(info.PermanentInstances)
	transientInstances := copyInstances(info.TransientInstances)
	logger.Info(
		"copy of instances made",
		zap.Int("permanentInstancesCount", len(permanentInstances)),
		zap.Int("transientInstancesCount", len(transientInstances)),
	)

	if !options.ConsiderFreeInstances {
		permanentInstances = removeFreeInstances(permanentInstances)
		transientInstances = removeFreeInstances(transientInstances)
		logger.Info(
			"removed free instances",
			zap.Int("remainingPermanentInstances", len(permanentInstances)),
			zap.Int("remainingTransientInstances", len(transientInstances)),
		)
	}

	allInstances := append(permanentInstances, transientInstances...)
	logger.Debug(
		"joined transient and permanent instances",
		zap.Int("totalInstanceCount", len(allInstances)),
	)

	advice := &schema.RegionAdvice{}

	for _, svc := range services {
		permanentCount := svc.MinInstances
		transientCount := svc.MaxInstances - svc.MinInstances

		logger.Info(
			"advising for service",
			zap.Any("service", svc),
			zap.Int("permanentInstanceCount", permanentCount),
			zap.Int("transientInstanceCount", transientCount),
		)

		for i := 0; i < permanentCount; i += 1 {
			logger.Info(
				"selecting permanent instance for service",
				zap.String("serviceName", svc.Name),
			)

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
			logger.Info(
				"selected and assigned permanent instance for service",
				zap.String("serviceName", svc.Name),
				zap.Any("instance", selectedInstance),
			)

			if options.ShareInstancesBetweenServices {
				sharedInstance := createSharedInstance(selectedInstance, svc)
				permanentInstances = append(permanentInstances, sharedInstance)
				logger.Info(
					"added shared instance to available instances",
					zap.String("originalInstanceId", selectedInstance.Id),
					zap.Any("newSharedInstance", sharedInstance),
				)
			}
		}

		for i := 0; i < transientCount; i += 1 {
			logger.Info(
				"selecting transient instance for service",
				zap.String("serviceName", svc.Name),
			)

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
			logger.Info(
				"selected and assigned permanent instance for service",
				zap.String("serviceName", svc.Name),
				zap.Any("instance", selectedInstance),
			)

			if options.AvoidRepeatedInstanceTypes {
				allInstances = removeInstancesWithName(allInstances, selectedInstance.Name)
				logger.Info(
					"removed instances of same type from available transient instances",
					zap.String("instanceId", selectedInstance.Id),
					zap.String("nameRemoved", selectedInstance.Name),
				)
			}

			if options.ShareInstancesBetweenServices {
				sharedInstance := createSharedInstance(selectedInstance, svc)

				allInstances = append(allInstances, sharedInstance)
				if isPermanentInstance(selectedInstance) {
					permanentInstances = append(
						permanentInstances,
						createSharedInstance(selectedInstance, svc),
					)
				}

				logger.Info(
					"added shared instance to available instances",
					zap.String("originalInstanceId", selectedInstance.Id),
					zap.Any("newSharedInstance", sharedInstance),
				)
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

func removeInstancesWithName(
	instances []*instPkg.Instance,
	name string,
) []*instPkg.Instance {
	filtered := []*instPkg.Instance{}
	for _, inst := range instances {
		if inst.Name != name {
			filtered = append(filtered, inst)
		}
	}
	return filtered
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
	logger *zap.Logger,
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
