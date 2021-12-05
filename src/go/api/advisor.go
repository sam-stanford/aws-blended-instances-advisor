package api

import "fmt"

type Advisor struct {
	Type        AdvisorType `json:"type"`
	Focus       string      `json:"focus"`
	FocusWeight float64     `json:"focusWeight"`
}

type AdvisorType string

const (
	Random   AdvisorType = "random"
	Weighted AdvisorType = "weighted"
	// TODO: "Focus" & "Custom" w/ custom configs
)

// TODO: Doc, test & complete
func (a *Advisor) Validate() error {
	if a.FocusWeight < 0 || a.FocusWeight > 1 {
		return fmt.Errorf(
			"advisor has focusWeight value outside of range of [0,1]: %f",
			a.FocusWeight,
		)
	}
	return nil
}
