package instances

import (
	"aws-blended-instances-advisor/api/schema"
	awsTypes "aws-blended-instances-advisor/aws/types"
)

// An Instance represents a single AWS EC2 instance offering.
type Instance struct {
	Id                    string          `json:"id"`
	Name                  string          `json:"name"`
	MemoryGb              float64         `json:"memory"`
	Vcpu                  int             `json:"vcpu"`
	Region                awsTypes.Region `json:"region"`
	AvailabilityZone      string          `json:"az"`
	OperatingSystem       string          `json:"os"`
	PricePerHour          float64         `json:"price"`
	RevocationProbability float64         `json:"revocProb"`
}

// ToApiSchemaInstance converts an Instance to an Instance suitable
// for the defined API specification.
func (inst *Instance) ToApiSchemaInstance() *schema.Instance {
	return &schema.Instance{
		Id:                    inst.Id,
		Name:                  inst.Name,
		MemoryGb:              inst.MemoryGb,
		Vcpu:                  inst.Vcpu,
		Region:                inst.Region.CodeString(),
		AvailabilityZone:      inst.AvailabilityZone,
		OperatingSystem:       inst.OperatingSystem,
		PricePerHour:          inst.PricePerHour,
		RevocationProbability: inst.RevocationProbability,
	}
}

// MakeCopy returns a copy of the Instance.
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
		RevocationProbability: inst.RevocationProbability,
	}
}
