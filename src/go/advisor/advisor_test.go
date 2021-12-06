package advisor

import "testing"

func TestWeightedIsAdvisor(t *testing.T) {
	var _ Advisor = WeightedAdvisor{}
}
