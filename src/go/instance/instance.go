package instance

import (
	awsTypes "ec2-test/aws/types"
)

type RegionInstancesMap map[awsTypes.Region][]Instance // TODO: Use

// TODO: Doc
type Instance struct {
	Name                  string                   `json:"name"`
	MemoryGb              float64                  `json:"memory"`
	Vcpu                  int                      `json:"vcpu"`
	Region                awsTypes.Region          `json:"region"`
	AvailabilityZone      string                   `json:"az"`
	OperatingSystem       awsTypes.OperatingSystem `json:"os"`
	PricePerHour          float64                  `json:"price"`
	RevocationProbability float64                  `json:"revocProb"`
}
