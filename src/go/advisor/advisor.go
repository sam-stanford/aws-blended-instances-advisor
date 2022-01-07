package advisor

import (
	"aws-blended-instances-advisor/api/schema"
	instPkg "aws-blended-instances-advisor/instances"

	"go.uber.org/zap"
)

// An Advisor is used to select one or more Instances from a group of Instances,
// and can also be used to score selections.
type Advisor interface {
	// Advise selects and scores Instances from a group of available
	// Instances for all Regions, returning the selection and information as an
	// Advice.
	Advise(
		instancesInfo instPkg.GlobalInfo,
		services []schema.Service,
		options schema.Options,
		logger *zap.Logger,
	) (
		*schema.Advice,
		error,
	)

	// AdviseForRegion selects and scores Instances from a group of available
	// Instances for one Region, returning the selection and information as a
	// RegionAdvice.
	AdviseForRegion(
		info instPkg.RegionInfo,
		services []schema.Service,
		options schema.Options,
		logger *zap.Logger,
	) (
		*schema.RegionAdvice,
		error,
	)

	// ScoreRegionAdvice scores a selection of Instances (as a RegionAdvice),
	// returning an arbitrary score.
	//
	// The returned score can be used to compare RegionAdvices, with higher scores
	// meaning a better selection.
	ScoreRegionAdvice(
		advice *schema.RegionAdvice,
		globalAgg instPkg.Aggregates,
		services []schema.Service,
		logger *zap.Logger,
	) float64
}

// New creates an advisor, using the type provided in the info argument
// to determine which advisor to use.
func New(info schema.Advisor) Advisor {
	switch info.Type {
	case schema.Weighted:
		return NewWeightedAdvisor(info.Weights)

	default:
		return NewWeightedAdvisor(info.Weights)
	}
}
