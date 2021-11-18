package api

import (
	"ec2-test/aws/types"
	"testing"
)

const (
	CONFIG_FILEPATH          = "../../../../config.json"
	CONFIG_API_DOWNLOADS_DIR = "../../../../assets/downloads/"
	MAX_INSTANCES            = 100
	REGION1                  = "us-west-1"
	REGION2                  = "ca-central-1" // TODO: Why no on-demand for "eu-west-2"
)

func getTestRegions() (*types.Region, *types.Region, error) {
	r1, err := types.NewRegion(REGION1)
	if err != nil {
		return nil, nil, err
	}
	r2, err := types.NewRegion(REGION2)
	if err != nil {
		return nil, nil, err
	}
	return &r1, &r2, nil
}

func TestGetInstances(t *testing.T) {
	// TODO
}
