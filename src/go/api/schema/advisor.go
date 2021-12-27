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
	Random   AdvisorType = "random"
	Weighted AdvisorType = "weighted"
)

// TODO: Doc & test
func (a *Advisor) Validate() error {
	return nil // Nothing to validate
}
