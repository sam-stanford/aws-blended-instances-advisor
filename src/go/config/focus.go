package config

import "fmt"

// TODO: Doc
type ServiceFocus int

const (
	Availability ServiceFocus = iota
	Balanced
	Cost
	Performance
)

func (sf ServiceFocus) String() string {
	switch sf {
	case Availability:
		return "availability"
	case Balanced:
		return "balanced"
	case Cost:
		return "cost"
	case Performance:
		return "performance"
	default:
		return "balanced"
	}
}

func ParseServiceFocus(serviceFocusString string) (ServiceFocus, error) {
	switch serviceFocusString {
	case "availability":
		return Availability, nil
	case "balanced":
		return Balanced, nil
	case "cost":
		return Cost, nil
	case "performance":
		return Performance, nil
	default:
		return -1, fmt.Errorf("\"%s\" is not a ServiceFocus", serviceFocusString)
	}
}
