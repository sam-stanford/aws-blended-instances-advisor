package types

import (
	"errors"
	"fmt"
)

type OperatingSystem int

const (
	Linux OperatingSystem = iota
	Windows
)

func (os OperatingSystem) ToString() string {
	switch os {
	case Linux:
		return "Linux"
	case Windows:
		return "Windows"
	default:
		return "NO_OS"
	}
}

func NewOperatingSystemFromString(value string) (OperatingSystem, error) {
	switch value {
	case "Linux":
		return Linux, nil
	case "Windows":
		return Windows, nil
	}

	return -1, errors.New(
		fmt.Sprintf("provided value of %s does not match any operating system", value))
}
