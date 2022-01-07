package schema

import (
	awsTypes "aws-blended-instances-advisor/aws/types"
)

type Options struct {
	AvoidRepeatedInstanceTypes    bool     `json:"avoidRepeatedInstanceTypes"`
	ShareInstancesBetweenServices bool     `json:"shareInstancesBetweenServices"`
	ConsiderFreeInstances         bool     `json:"considerFreeInstances"`
	Regions                       []string `json:"regions"`
}

// Validate checks that an Options variable is well-formed
// and is true to the API specification.
func (o *Options) Validate() error {
	_, err := awsTypes.NewRegions(o.Regions)
	return err
}

// GetRegionsAsAwsRegions converts an Options' regions strings into
// Regions used in the aws/api package.
func (o *Options) GetRegionsAsAwsRegions() ([]awsTypes.Region, error) {
	return awsTypes.NewRegions(o.Regions)
}
