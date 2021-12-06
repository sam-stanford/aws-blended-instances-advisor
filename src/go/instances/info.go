package instances

import awsTypes "ec2-test/aws/types"

// TODO: Docs

// TODO: Is this used?

type RegionInfoMap map[awsTypes.Region]RegionInfo

type RegionInfo struct {
	PermanentInstances InstancesAndAggregates `json:"permament"`
	AllInstances       InstancesAndAggregates `json:"all"`
}

type InstancesAndAggregates struct {
	Instances  []*Instance `json:"instances"`
	Aggregates Aggregates  `json:"aggregates"`
}

func CreateRegionInfo(permanentInstances []*Instance, transientInstances []*Instance) RegionInfo {
	allInstances := append(permanentInstances, transientInstances...)
	onDemandAggs := CalculateAggregates(permanentInstances)
	allInstancesAggs := CalculateAggregates(allInstances)

	return RegionInfo{
		PermanentInstances: InstancesAndAggregates{
			Instances:  permanentInstances,
			Aggregates: onDemandAggs,
		},
		AllInstances: InstancesAndAggregates{
			Instances:  allInstances,
			Aggregates: allInstancesAggs,
		},
	}
}
