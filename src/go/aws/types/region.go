package types

import "fmt"

// A Region represents one AWS region.
type Region int

const (
	UsEast1 Region = iota
	UsEast2
	UsWest1
	UsWest2
	ApSouth1
	ApNorthEast3
	ApNorthEast2
	ApSouthEast1
	ApSouthEast2
	ApNorthEast1
	CaCentral1
	EuCentral1
	EuWest1
	EuWest2
	EuWest3
	EuNorth1
	SaEast1
)

// GetAllRegions returns a slice containing all possible Regions.
func GetAllRegions() []Region {
	lastRegion := SaEast1

	r := []Region{}
	for i := 0; i <= int(lastRegion); i += 1 {
		r = append(r, Region(i))
	}
	return r
}

// CodeString returns the string representation of a Region
// in its code format.
//
// Example: eu-east-2
func (region Region) CodeString() string {
	switch region {
	case UsEast1:
		return "us-east-1"
	case UsEast2:
		return "us-east-2"
	case UsWest1:
		return "us-west-1"
	case UsWest2:
		return "us-west-2"
	case ApSouth1:
		return "ap-south-1"
	case ApNorthEast3:
		return "ap-northeast-3"
	case ApNorthEast2:
		return "ap-northeast-2"
	case ApNorthEast1:
		return "ap-southeast-1"
	case ApSouthEast2:
		return "ap-southeast-2"
	case ApSouthEast1:
		return "ap-northeast-1"
	case CaCentral1:
		return "ca-central-1"
	case EuCentral1:
		return "eu-central-1"
	case EuWest1:
		return "eu-west-1"
	case EuWest2:
		return "eu-west-2"
	case EuWest3:
		return "eu-west-3"
	case EuNorth1:
		return "eu-north-1"
	case SaEast1:
		return "sa-east-1"
	default:
		return "NO_REGION"
	}
}

// CodeString returns the string representation of a Region
// which is representative of the Region's name
//
// Example: US East (N. Virginia)
func (region Region) NameString() string {
	switch region {
	case UsEast1:
		return "US East (N. Virginia)"
	case UsEast2:
		return "US East (Ohio)"
	case UsWest1:
		return "US West (N. California)"
	case UsWest2:
		return "US West (Oregon)"
	case ApSouth1:
		return "Asia Pacific (Mumbai)"
	case ApNorthEast3:
		return "Asia Pacific (Osaka)"
	case ApNorthEast2:
		return "Asia Pacific (Seoul)"
	case ApNorthEast1:
		return "Asia Pacific (Singapore)"
	case ApSouthEast2:
		return "Asia Pacific (Sydney)"
	case ApSouthEast1:
		return "Asia Pacific (Tokyo)"
	case CaCentral1:
		return "Canada (Central)"
	case EuCentral1:
		return "EU (Frankfurt)"
	case EuWest1:
		return "EU (Ireland)"
	case EuWest2:
		return "EU (London)"
	case EuWest3:
		return "EU (Paris)"
	case EuNorth1:
		return "EU (Stockholm)"
	case SaEast1:
		return "South America (Sao Paulo)"
	default:
		return "NO_REGION"
	}
}

// NewRegion creates a new region from a string representation
// of the region, including name and code strings.
//
// An error is returned if the string does not match any region.
func NewRegion(value string) (Region, error) {
	switch value {
	case "us-east-1", "US East (N. Virginia)":
		return UsEast1, nil
	case "us-east-2", "US East (Ohio)":
		return UsEast2, nil
	case "us-west-1", "US West (N. California)":
		return UsWest1, nil
	case "us-west-2", "US West (Oregon)":
		return UsWest2, nil
	case "ap-south-1", "Asia Pacific (Mumbai)":
		return ApSouth1, nil
	case "ap-northeast-3", "Asia Pacific (Osaka)":
		return ApNorthEast3, nil
	case "ap-northeast-2", "Asia Pacific (Seoul)":
		return ApNorthEast2, nil
	case "ap-southeast-1", "Asia Pacific (Singapore)":
		return ApNorthEast1, nil
	case "ap-southeast-2", "Asia Pacific (Sydney)":
		return ApSouthEast2, nil
	case "ap-northeast-1", "Asia Pacific (Tokyo)":
		return ApSouthEast1, nil
	case "ca-central-1", "Canada (Central)":
		return CaCentral1, nil
	case "eu-central-1", "Europe (Frankfurt)", "EU (Frankfurt)":
		return EuCentral1, nil
	case "eu-west-1", "Europe (Ireland)", "EU (Ireland)":
		return EuWest1, nil
	case "eu-west-2", "Europe (London)", "EU (London)":
		return EuWest2, nil
	case "eu-west-3", "Europe (Paris)", "EU (Paris)":
		return EuWest3, nil
	case "eu-north-1", "Europe (Stockholm)", "EU (Stockholm)":
		return EuNorth1, nil
	case "sa-east-1", "South America (Sao Paulo)":
		return SaEast1, nil
	}

	return -1, fmt.Errorf("provided value of \"%s\" does not match any region", value)
}

// NewRegions creates multiple Regions from a slice of region string values.
//
// An error is returned if any value does not match a Region.
func NewRegions(values []string) ([]Region, error) {
	regions := []Region{}
	for _, value := range values {
		r, err := NewRegion(value)
		if err != nil {
			return nil, err
		}
		regions = append(regions, r)
	}
	return regions, nil
}
