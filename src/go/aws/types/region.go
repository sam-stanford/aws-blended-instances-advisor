package types

import "errors"

type Region int

const (
	UsEast1 Region = iota
	UsEast2
	UsWest1
	UsWest2
	AfSouth1
	ApEast1
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
	EuSouth1
	EuWest3
	EuNorth1
	MeSouth1
	SaEast1
)

func (region Region) ToCodeString() string {
	switch region {
	case UsEast1:
		return "us-east-1"
	case UsEast2:
		return "us-east-2"
	case UsWest1:
		return "us-west-1"
	case UsWest2:
		return "us-west-2"
	case AfSouth1:
		return "af-south-1"
	case ApEast1:
		return "ap-east-1"
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
	case EuSouth1:
		return "eu-south-1"
	case EuNorth1:
		return "eu-north-1"
	case MeSouth1:
		return "me-south-1"
	case SaEast1:
		return "sa-east-1"
	default:
		return "NO_REGION"
	}
}

func (region Region) ToNameString() string {
	switch region {
	case UsEast1:
		return "US East (N. Virginia)"
	case UsEast2:
		return "US East (Ohio)"
	case UsWest1:
		return "US West (N. California)"
	case UsWest2:
		return "US West (Oregon)"
	case AfSouth1:
		return "Africa (Cape Town)"
	case ApEast1:
		return "Asia Pacific (Hong Kong)"
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
		return "Europe (Frankfurt)"
	case EuWest1:
		return "Europe (Ireland)"
	case EuWest2:
		return "Europe (London)"
	case EuWest3:
		return "Europe (Paris)"
	case EuSouth1:
		return "Europe (Milan)"
	case EuNorth1:
		return "Europe (Stockholm)"
	case MeSouth1:
		return "Middle East (Bahrain)"
	case SaEast1:
		return "South America (São Paulo)"
	default:
		return "NO_REGION"
	}
}

func NewRegionFromString(value string) (Region, error) {
	switch value {
	case "us-east-1", "US East (N. Virginia)":
		return UsEast1, nil
	case "us-east-2", "US East (Ohio)":
		return UsEast2, nil
	case "us-west-1", "US West (N. California)":
		return UsWest1, nil
	case "us-west-2", "US West (Oregon)":
		return UsWest2, nil
	case "af-south-1", "Africa (Cape Town)":
		return AfSouth1, nil
	case "ap-east-1", "Asia Pacific (Hong Kong)":
		return ApEast1, nil
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
	case "eu-central-1", "Europe (Frankfurt)":
		return EuCentral1, nil
	case "eu-west-1", "Europe (Ireland)":
		return EuWest1, nil
	case "eu-west-2", "Europe (London)":
		return EuWest2, nil
	case "eu-west-3", "Europe (Paris)":
		return EuWest3, nil
	case "eu-south-1", "Europe (Milan)":
		return EuSouth1, nil
	case "eu-north-1", "Europe (Stockholm)":
		return EuNorth1, nil
	case "me-south-1", "Middle East (Bahrain)":
		return MeSouth1, nil
	case "sa-east-1", "South America (São Paulo)":
		return SaEast1, nil
	}

	return -1, errors.New("provided value does not match any region")
}
