package api

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
