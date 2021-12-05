package instances

import (
	"ec2-test/api"
	awsTypes "ec2-test/aws/types"
)

// TODO: Doc
type Instance struct {
	Id                    string                   `json:"id"`
	Name                  string                   `json:"name"`
	MemoryGb              float64                  `json:"memory"`
	Vcpu                  int                      `json:"vcpu"`
	Region                awsTypes.Region          `json:"region"`
	AvailabilityZone      string                   `json:"az"`
	OperatingSystem       awsTypes.OperatingSystem `json:"os"`
	PricePerHour          float64                  `json:"price"`
	RevocationProbability float64                  `json:"revocProb"`
}

// TODO: Doc
func (inst *Instance) ToApiInstance() *api.Instance {
	return &api.Instance{
		Id:                    inst.Id,
		Name:                  inst.Name,
		MemoryGb:              inst.MemoryGb,
		Vcpu:                  inst.Vcpu,
		Region:                inst.Region.ToCodeString(),
		AvailabilityZone:      inst.AvailabilityZone,
		OperatingSystem:       inst.OperatingSystem.ToString(),
		PricePerHour:          inst.PricePerHour,
		RevocationProbability: inst.RevocationProbability,
	}
}
