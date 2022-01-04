package utils

import "testing"

type floatsEqualTest struct {
	a     float64
	b     float64
	equal bool
}

func TestFloatsEqual(t *testing.T) {
	tests := []floatsEqualTest{
		{
			a:     4.0,
			b:     4.0,
			equal: true,
		},
		{
			a:     12.0 / 3.0,
			b:     16.0 / 4.0,
			equal: true,
		},
		{
			a:     4.0,
			b:     4.001,
			equal: false,
		},
	}

	for _, test := range tests {
		if FloatsEqual(test.a, test.b) != test.equal {
			if test.equal {
				t.Fatalf("Floats not considered equal: %f, %f", test.a, test.b)
			} else {
				t.Fatalf("Floats considered equal: %f, %f", test.a, test.b)
			}
		}
	}
}

type minMaxIntTest struct {
	a      int
	b      int
	wanted int
}

func TestMinOfInts(t *testing.T) {
	tests := []minMaxIntTest{
		{
			a:      0,
			b:      1,
			wanted: 0,
		},
		{
			a:      4,
			b:      3,
			wanted: 3,
		},
		{
			a:      9,
			b:      9,
			wanted: 9,
		},
		{
			a:      -2,
			b:      -9,
			wanted: -9,
		},
	}

	for _, test := range tests {
		got := MinOfInts(test.a, test.b)
		if got != test.wanted {
			t.Fatalf("Non-minimum value returned. Wanted: %d, got: %d", test.wanted, got)
		}
	}
}

func TestMaxOfInts(t *testing.T) {
	tests := []minMaxIntTest{
		{
			a:      0,
			b:      1,
			wanted: 1,
		},
		{
			a:      4,
			b:      3,
			wanted: 4,
		},
		{
			a:      9,
			b:      9,
			wanted: 9,
		},
		{
			a:      -2,
			b:      -9,
			wanted: -2,
		},
	}

	for _, test := range tests {
		got := MaxOfInts(test.a, test.b)
		if got != test.wanted {
			t.Fatalf("Non-minimum value returned. Wanted: %d, got: %d", test.wanted, got)
		}
	}
}

type minMaxFloatTest struct {
	a      float64
	b      float64
	wanted float64
}

func TestMinOfFloats(t *testing.T) {
	tests := []minMaxFloatTest{
		{
			a:      0.0,
			b:      0.8,
			wanted: 0,
		},
		{
			a:      4.5,
			b:      2.9,
			wanted: 2.9,
		},
		{
			a:      9.01,
			b:      9.01,
			wanted: 9.01,
		},
		{
			a:      -2.0,
			b:      -9.01,
			wanted: -9.01,
		},
	}

	for _, test := range tests {
		got := MinOfFloats(test.a, test.b)
		if got != test.wanted {
			t.Fatalf("Non-minimum value returned. Wanted: %f, got: %f", test.wanted, got)
		}
	}
}

func TestMaxOfFloats(t *testing.T) {
	tests := []minMaxFloatTest{
		{
			a:      0.0,
			b:      0.8,
			wanted: 0.8,
		},
		{
			a:      4.5,
			b:      2.9,
			wanted: 4.5,
		},
		{
			a:      9.01,
			b:      9.01,
			wanted: 9.01,
		},
		{
			a:      -2.0,
			b:      -9.01,
			wanted: -2.0,
		},
	}

	for _, test := range tests {
		got := MaxOfFloats(test.a, test.b)
		if got != test.wanted {
			t.Fatalf("Non-minimum value returned. Wanted: %f, got: %f", test.wanted, got)
		}
	}
}
