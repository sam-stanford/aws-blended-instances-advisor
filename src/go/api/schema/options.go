package schema

import (
	awsTypes "aws-blended-instances-advisor/aws/types"
)

type Options struct {
	AvoidRepeatedInstanceTypes       bool     `json:"avoidRepeatedInstanceTypes"`
	ShareInstancesBetweenServices    bool     `json:"shareInstancesBetweenServices"`
	ConsiderFreeInstances            bool     `json:"considerFreeInstances"`
	ShareInstancesBetweenSameService bool     `json:"shareInstancesBetweenSameService"` // TODO: Implement
	Regions                          []string `json:"regions"`
}

// TODO: Test & doc

func (o *Options) Validate() error {
	_, err := awsTypes.NewRegions(o.Regions)
	return err
}

func (o *Options) GetRegionsAsAwsRegions() ([]awsTypes.Region, error) {
	return awsTypes.NewRegions(o.Regions)
}
