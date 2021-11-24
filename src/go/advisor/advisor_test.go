package advisor

import "testing"

func TestNaiveReliabilityAdvisorIsAdvisor(t *testing.T) {
	var _ Advisor = NaiveReliabilityAdvisor{}
}
