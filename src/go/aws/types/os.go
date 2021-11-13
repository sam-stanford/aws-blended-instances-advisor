package types

import (
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
	case "Linux", "RHEL", "Red Hat Enterprise Linux with HA", "SUSE":
		return Linux, nil
	case "Windows":
		return Windows, nil
	}

	return -1, fmt.Errorf("provided value of %s does not match any operating system", value)
}
