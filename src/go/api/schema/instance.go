package schema

type Instance struct {
	Id                    string  `json:"id"` // TODO: Add to instances pkg & generate on creation
	Name                  string  `json:"name"`
	MemoryGb              float64 `json:"memory"`
	Vcpu                  int     `json:"vcpu"`
	Region                string  `json:"region"`
	AvailabilityZone      string  `json:"az"`
	OperatingSystem       string  `json:"os"`
	PricePerHour          float64 `json:"price"`
	RevocationProbability float64 `json:"revocProb"`
}
