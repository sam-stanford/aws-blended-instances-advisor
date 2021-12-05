package instances

import (
	"ec2-test/utils"
	"testing"
)

type aggregatesTest struct {
	instances []*Instance
	expected  Aggregates
}

type normaliseTest struct {
	aggregates Aggregates
	instance   Instance
	expected   float64
}

func TestCalculateAggregates(t *testing.T) {
	tests := map[string]aggregatesTest{
		"single instance": {
			instances: []*Instance{
				{Vcpu: 4, RevocationProbability: 0, PricePerHour: 0.001},
			},
			expected: Aggregates{
				Count: 1,

				MinVcpu:  4,
				MaxVcpu:  4,
				MeanVcpu: 4,

				MinRevocationProbability:  0,
				MaxRevocationProbability:  0,
				MeanRevocationProbability: 0,

				MinPricePerHour:  0.001,
				MaxPricePerHour:  0.001,
				MeanPricePerHour: 0.001,
			},
		},
		"multiple instances": {
			instances: []*Instance{
				{Vcpu: 4, RevocationProbability: 0, PricePerHour: 0.001},
				{Vcpu: 4, RevocationProbability: 0.1, PricePerHour: 0.005},
				{Vcpu: 8, RevocationProbability: 0.2, PricePerHour: 0.01},
				{Vcpu: 16, RevocationProbability: 0.3, PricePerHour: 0.05},
			},
			expected: Aggregates{
				Count: 4,

				MinVcpu:  4,
				MaxVcpu:  16,
				MeanVcpu: 8,

				MinRevocationProbability:  0,
				MaxRevocationProbability:  0.3,
				MeanRevocationProbability: 0.15,

				MinPricePerHour:  0.001,
				MaxPricePerHour:  0.05,
				MeanPricePerHour: 0.0165,
			},
		},
	}

	for name, test := range tests {
		agg := CalculateAggregates(test.instances)

		if agg.Count != test.expected.Count {
			t.Fatalf(
				"Aggregate count not equal for test \"%s\".Wanted: %d, got: %d",
				name,
				test.expected.Count,
				agg.Count,
			)
		}

		if agg.MinVcpu != test.expected.MinVcpu {
			t.Fatalf(
				"Aggregate min VCPU value not equal for test \"%s\".Wanted: %d, got: %d",
				name,
				test.expected.MinVcpu,
				agg.MinVcpu,
			)
		}

		if agg.MaxVcpu != test.expected.MaxVcpu {
			t.Fatalf(
				"Aggregate max VCPU value not equal for test \"%s\".Wanted: %d, got: %d",
				name,
				test.expected.MaxVcpu,
				agg.MaxVcpu,
			)
		}

		if !utils.FloatsEqual(agg.MeanVcpu, test.expected.MeanVcpu) {
			t.Fatalf(
				"Aggregate mean VCPU value not equal for test \"%s\".Wanted: %f, got: %f",
				name,
				test.expected.MeanVcpu,
				agg.MeanVcpu,
			)
		}

		if !utils.FloatsEqual(agg.MinRevocationProbability, test.expected.MinRevocationProbability) {
			t.Fatalf(
				"Aggregate min revocation probability not equal for test \"%s\".Wanted: %f, got: %f",
				name,
				test.expected.MinRevocationProbability,
				agg.MinRevocationProbability,
			)
		}

		if !utils.FloatsEqual(agg.MaxRevocationProbability, test.expected.MaxRevocationProbability) {
			t.Fatalf(
				"Aggregate max revocation probability not equal for test \"%s\".Wanted: %f, got: %f",
				name,
				test.expected.MaxRevocationProbability,
				agg.MaxRevocationProbability,
			)
		}

		if !utils.FloatsEqual(agg.MeanRevocationProbability, test.expected.MeanRevocationProbability) {
			t.Fatalf(
				"Aggregate mean revocation probability not equal for test \"%s\".Wanted: %f, got: %f",
				name,
				test.expected.MeanRevocationProbability,
				agg.MeanRevocationProbability,
			)
		}

		if !utils.FloatsEqual(agg.MinPricePerHour, test.expected.MinPricePerHour) {
			t.Fatalf(
				"Aggregate min price per hour not equal for test \"%s\".Wanted: %f, got: %f",
				name,
				test.expected.MinPricePerHour,
				agg.MinPricePerHour,
			)
		}

		if !utils.FloatsEqual(agg.MaxPricePerHour, test.expected.MaxPricePerHour) {
			t.Fatalf(
				"Aggregate max price per hour not equal for test \"%s\".Wanted: %f, got: %f",
				name,
				test.expected.MaxPricePerHour,
				agg.MaxPricePerHour,
			)
		}

		if !utils.FloatsEqual(agg.MeanPricePerHour, test.expected.MeanPricePerHour) {
			t.Fatalf(
				"Aggregate mean price per hour not equal for test \"%s\".Wanted: %f, got: %f",
				name,
				test.expected.MeanPricePerHour,
				agg.MeanPricePerHour,
			)
		}
	}
}

func TestNormaliseVcpu(t *testing.T) {
	tests := map[string]normaliseTest{
		"instance equals all aggregates": {
			aggregates: Aggregates{Count: 1, MinVcpu: 4, MaxVcpu: 4, MeanVcpu: 4}, // Values to form: 4
			instance:   Instance{Vcpu: 4},
			expected:   1, // 1 / count
		},
		"instance is min of aggregates": {
			aggregates: Aggregates{Count: 3, MinVcpu: 4, MaxVcpu: 12, MeanVcpu: 8}, // Values to form: 4, 8, 12
			instance:   Instance{Vcpu: 4},
			expected:   0,
		},
		"instance is max of aggregates": {
			aggregates: Aggregates{Count: 3, MinVcpu: 4, MaxVcpu: 12, MeanVcpu: 8}, // Values to form: 4, 8, 12
			instance:   Instance{Vcpu: 12},
			expected:   1,
		},
		"instance is middle of aggregates": {
			aggregates: Aggregates{Count: 3, MinVcpu: 4, MaxVcpu: 12, MeanVcpu: 8}, // Values to form: 4, 8, 12
			instance:   Instance{Vcpu: 8},
			expected:   0.5,
		},
	}

	for name, test := range tests {
		got := test.aggregates.NormaliseVcpu(test.instance.Vcpu)
		if !utils.FloatsEqual(got, test.expected) {
			t.Fatalf(
				"Normalised value is incorrect for test \"%s\". Wanted: %f, got: %f",
				name,
				test.expected,
				got,
			)
		}
	}
}

func TestNormaliseRevocationProbability(t *testing.T) {
	tests := map[string]normaliseTest{
		"instance equals all aggregates": {
			aggregates: Aggregates{Count: 1, MinRevocationProbability: 0.1, MaxRevocationProbability: 0.1, MeanRevocationProbability: 0.1}, // Values to form: 0.1
			instance:   Instance{RevocationProbability: 4},
			expected:   1, // 1 / count
		},
		"instance is min of aggregates": {
			aggregates: Aggregates{Count: 3, MinRevocationProbability: 0.1, MaxRevocationProbability: 0.3, MeanRevocationProbability: 0.2}, // Values to form: 0.1, 0.2, 0.3
			instance:   Instance{RevocationProbability: 0.1},
			expected:   0,
		},
		"instance is max of aggregates": {
			aggregates: Aggregates{Count: 3, MinRevocationProbability: 0.1, MaxRevocationProbability: 0.3, MeanRevocationProbability: 0.2}, // Values to form: 0.1, 0.2, 0.3
			instance:   Instance{RevocationProbability: 0.3},
			expected:   1,
		},
		"instance is middle of aggregates": {
			aggregates: Aggregates{Count: 3, MinRevocationProbability: 0.1, MaxRevocationProbability: 0.3, MeanRevocationProbability: 0.2}, // Values to form: 0.1, 0.2, 0.3
			instance:   Instance{RevocationProbability: 0.2},
			expected:   0.5,
		},
	}

	for name, test := range tests {
		got := test.aggregates.NormaliseRevocationProbability(test.instance.RevocationProbability)
		if !utils.FloatsEqual(got, test.expected) {
			t.Fatalf(
				"Normalised value is incorrect for test \"%s\". Wanted: %f, got: %f",
				name,
				test.expected,
				got,
			)
		}
	}
}

func TestNormalisePricePerHour(t *testing.T) {
	tests := map[string]normaliseTest{
		"instance equals all aggregates": {
			aggregates: Aggregates{Count: 1, MinPricePerHour: 0.1, MaxPricePerHour: 0.1, MeanPricePerHour: 0.1}, // Values to form: 0.1
			instance:   Instance{PricePerHour: 4},
			expected:   1, // 1 / count
		},
		"instance is min of aggregates": {
			aggregates: Aggregates{Count: 3, MinPricePerHour: 0.1, MaxPricePerHour: 0.3, MeanPricePerHour: 0.2}, // Values to form: 0.1, 0.2, 0.3
			instance:   Instance{PricePerHour: 0.1},
			expected:   0,
		},
		"instance is max of aggregates": {
			aggregates: Aggregates{Count: 3, MinPricePerHour: 0.1, MaxPricePerHour: 0.3, MeanPricePerHour: 0.2}, // Values to form: 0.1, 0.2, 0.3
			instance:   Instance{PricePerHour: 0.3},
			expected:   1,
		},
		"instance is middle of aggregates": {
			aggregates: Aggregates{Count: 3, MinPricePerHour: 0.1, MaxPricePerHour: 0.3, MeanPricePerHour: 0.2}, // Values to form: 0.1, 0.2, 0.3
			instance:   Instance{PricePerHour: 0.2},
			expected:   0.5,
		},
	}

	for name, test := range tests {
		got := test.aggregates.NormalisePricePerHour(test.instance.PricePerHour)
		if !utils.FloatsEqual(got, test.expected) {
			t.Fatalf(
				"Normalised value is incorrect for test \"%s\". Wanted: %f, got: %f",
				name,
				test.expected,
				got,
			)
		}
	}
}
