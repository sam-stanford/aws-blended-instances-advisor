package instances

import (
	"aws-blended-instances-advisor/api/schema"
	awsTypes "aws-blended-instances-advisor/aws/types"
)

// TODO: Doc
type Instance struct {
	Id                    string                   `json:"id"`
	Name                  string                   `json:"name"`
	MemoryGb              float64                  `json:"memory"`
	Vcpu                  int                      `json:"vcpu"`
	Region                awsTypes.Region          `json:"region"`
	AvailabilityZone      string                   `json:"az"`
	OperatingSystem       awsTypes.OperatingSystem `json:"os"` // TODO: Change to string or provide search option
	PricePerHour          float64                  `json:"price"`
	RevocationProbability float64                  `json:"revocProb"`
}

// TODO: Doc
func (inst *Instance) ToApiSchemaInstance() *schema.Instance {
	return &schema.Instance{
		Id:                    inst.Id,
		Name:                  inst.Name,
		MemoryGb:              inst.MemoryGb,
		Vcpu:                  inst.Vcpu,
		Region:                inst.Region.CodeString(), // TODO: Rename to CodeString
		AvailabilityZone:      inst.AvailabilityZone,
		OperatingSystem:       inst.OperatingSystem.ToString(), // TODO: Rename to String
		PricePerHour:          inst.PricePerHour,
		RevocationProbability: inst.RevocationProbability,
	}
}

// TODO: Doc
func (inst *Instance) MakeCopy() *Instance {
	return &Instance{
		Id:                    inst.Id,
		Name:                  inst.Name,
		MemoryGb:              inst.MemoryGb,
		Vcpu:                  inst.Vcpu,
		Region:                inst.Region,
		AvailabilityZone:      inst.AvailabilityZone,
		OperatingSystem:       inst.OperatingSystem,
		PricePerHour:          inst.PricePerHour,
		RevocationProbability: inst.PricePerHour,
	}
}
