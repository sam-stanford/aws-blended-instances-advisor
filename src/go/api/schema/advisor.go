package schema

type Advisor struct {
	Type    AdvisorType `json:"type"`
	Weights AdvisorWeights
}

type AdvisorWeights struct {
	Price        float64 `json:"price"`
	Availability float64 `json:"availability"`
	Performance  float64 `json:"performance"`
}

type AdvisorType string

const (
	Weighted AdvisorType = "weighted"
)

// Validate checks that an Advisor is well-formed
// and is true to the API specification.
func (a *Advisor) Validate() error {
	return nil // Nothing to validate
}
