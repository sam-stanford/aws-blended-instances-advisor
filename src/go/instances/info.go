package instances

import (
	"aws-blended-instances-advisor/aws/types"
	awsTypes "aws-blended-instances-advisor/aws/types"
	"aws-blended-instances-advisor/utils"
	"fmt"

	"go.uber.org/zap"
)

// GlobalInfo describes the instances and aggregates for all regions' instances.
type GlobalInfo struct {
	RegionInfoMap    RegionInfoMap `json:"regionInfoMap"`
	GlobalAggregates Aggregates    `json:"globalAggregates"`
}

// RegionInfoMap is a map between Regions and their respective RegionInfo.
type RegionInfoMap map[awsTypes.Region]RegionInfo

// RegionInfo describes the instances and aggregate for a given region's instances.
type RegionInfo struct {
	PermanentInstances []*Instance `json:"permanentInstances"`
	TransientInstances []*Instance `json:"transientInstances"`

	PermanentAggregates Aggregates `json:"permanentAggregates"`
	TransientAggregates Aggregates `json:"transientAggregates"`
	RegionAggregates    Aggregates `json:"regionAggregates"`
}

// CreateGlobalInfo creates a new GlobalInfo.
func CreateGlobalInfo(
	regionToPersistentInstancesMap map[types.Region][]*Instance,
	regionToTransientInstancesMap map[types.Region][]*Instance,
	regions []types.Region,
) GlobalInfo {

	regionInfoMap := CreateRegionInfoMap(
		regionToPersistentInstancesMap,
		regionToTransientInstancesMap,
		regions,
	)

	return GlobalInfo{
		RegionInfoMap:    regionInfoMap,
		GlobalAggregates: CalculateGlobalAggregates(regionInfoMap),
	}
}

// CreateRegionInfoMap creates a new RegionInfoMap from
// Region-to-Instance maps for on-demand and spot instances.
func CreateRegionInfoMap(
	regionToPersistentInstancesMap map[types.Region][]*Instance,
	regionToTransientInstancesMap map[types.Region][]*Instance,
	regions []types.Region,
) RegionInfoMap {

	regionInfoMap := make(RegionInfoMap)
	for _, region := range regions {

		thisRegionInfo := CreateRegionInfo(
			regionToPersistentInstancesMap[region],
			regionToTransientInstancesMap[region],
		)
		regionInfoMap[region] = thisRegionInfo
	}
	return regionInfoMap
}

// CreateRegionInfo creates a new RegionInfo.
func CreateRegionInfo(permanentInstances []*Instance, transientInstances []*Instance) RegionInfo {
	permanentAggs := CalculateAggregates(permanentInstances)
	transientAggs := CalculateAggregates(transientInstances)
	regionAggs := CombineAggregates([]Aggregates{permanentAggs, transientAggs})

	return RegionInfo{
		PermanentInstances:  permanentInstances,
		TransientInstances:  transientInstances,
		PermanentAggregates: permanentAggs,
		TransientAggregates: transientAggs,
		RegionAggregates:    regionAggs,
	}
}

// CalculateGlobalAggregates calculates aggregates for all instances in a RegionInfoMap.
func CalculateGlobalAggregates(regionInfoMap RegionInfoMap) Aggregates {
	allAggs := []Aggregates{}
	for _, info := range regionInfoMap {
		allAggs = append(allAggs, info.RegionAggregates)
	}
	return CombineAggregates(allAggs)
}

// Log logs information regarding the GlobalInfo, including information relating
// to each region described by the GlobalInfo.
func (info *GlobalInfo) Log(message string, logger *zap.Logger) {
	logger.Info(
		message,
		zap.Int("regionCount", len(info.RegionInfoMap)),
		zap.Any("globalAggregates", info.GlobalAggregates),
	)
	for region, regionInfo := range info.RegionInfoMap {
		logger.Info(
			message,
			zap.String("region", region.CodeString()),
			zap.Int("transientInstanceCount", len(regionInfo.TransientInstances)),
			zap.Int("permanentInstanceCount", len(regionInfo.PermanentInstances)),
			zap.Any("transientAggregates", regionInfo.TransientAggregates),
			zap.Any("permanentAggregates", regionInfo.PermanentAggregates),
		)
	}
}

// TODO: Test
// Validate checks the GlobalInfo, ensuring the described aggregates
// correctly reflect the described instances.
//
// An error is returned if the GlobalInfo is invalid.
func (info *GlobalInfo) Validate() error {
	totalInstances := 0

	for region, regionInfo := range info.RegionInfoMap {
		err := regionInfo.Validate()
		if err != nil {
			return utils.PrependToError(
				err,
				fmt.Sprintf("info for region %s invalid", region.CodeString()),
			)
		}
		totalInstances += regionInfo.PermanentAggregates.Count + regionInfo.TransientAggregates.Count
	}

	if totalInstances != info.GlobalAggregates.Count {
		return fmt.Errorf(
			"different total instance count than described by aggregates (%d vs %d)",
			totalInstances,
			info.GlobalAggregates.Count,
		)
	}

	return nil
}

// TODO: Test
// Validate checks the RegionInfo, ensuring the aggregates correctly
// describe the listed instances.
//
// An error is returned if the RegionInfo is invalid.
func (info *RegionInfo) Validate() error {
	numPermanentInstances := len(info.PermanentInstances)
	numTransientInstances := len(info.TransientInstances)
	numRegionInstances := numPermanentInstances + numTransientInstances

	if numPermanentInstances == 0 {
		return fmt.Errorf("no permanent instances")
	}
	if numTransientInstances == 0 {
		return fmt.Errorf("no transient instances")
	}

	if numPermanentInstances != info.PermanentAggregates.Count {
		return fmt.Errorf(
			"different permanent instance count than described by aggregates (%d vs %d)",
			numPermanentInstances,
			info.PermanentAggregates.Count,
		)
	}
	if numTransientInstances != info.TransientAggregates.Count {
		return fmt.Errorf(
			"different transient instance count than described by aggregates (%d vs %d)",
			numPermanentInstances,
			info.PermanentAggregates.Count,
		)
	}
	if numRegionInstances != info.RegionAggregates.Count {
		return fmt.Errorf(
			"different region instance count than described by aggregates (%d vs %d)",
			numRegionInstances,
			info.RegionAggregates.Count,
		)
	}

	return nil
}
