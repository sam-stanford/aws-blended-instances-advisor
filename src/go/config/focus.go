package config

import "fmt"

// TODO: Docs
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

// TODO: Use this same setup for OS and Regions (i.e. from string and validate, rather than error every time)

func ServiceFocusFromString(serviceFocusString string) ServiceFocus {
	switch serviceFocusString {
	case "availability":
		return Availability
	case "cost":
		return Cost
	case "performance":
		return Performance
	default:
		return Balanced

	}
}

func ValidateServiceFocus(serviceFocusString string) error {
	switch serviceFocusString {
	case "availability", "balanced", "cost", "performance":
		return nil
	default:
		return fmt.Errorf("\"%s\" is not a ServiceFocus", serviceFocusString)
	}
}
