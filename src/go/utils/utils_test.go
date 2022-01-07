package utils

import (
	"errors"
	"testing"
)

func TestPrependToError(t *testing.T) {
	str1 := "prev err"
	str2 := "new err"
	joined := "new err: prev err"

	err1 := errors.New(str1)
	err2 := PrependToError(err1, str2)
	if err2.Error() != joined {
		t.Fatalf("String not prepended correctly. Wanted: %s, got: %s", joined, err2.Error())
	}
}

func TestCreateMockLogger(t *testing.T) {
	logger, err := CreateMockLogger()
	if err != nil {
		t.Fatalf("Error returned when creating logger: %s", err.Error())
	}
	logger.Info("") // Should not throw error
}

type validateIndexesTest struct {
	sliceLength, start, end int
	shouldError             bool
}

func TestValidateIndexes(t *testing.T) {
	tests := map[string]validateIndexesTest{
		"valid": {
			sliceLength: 4,
			start:       0,
			end:         2,
			shouldError: false,
		},
		"valid, full slice": {
			sliceLength: 2,
			start:       0,
			end:         2,
			shouldError: false,
		},
		"subslice length 0": {
			sliceLength: 1,
			start:       1,
			end:         1,
			shouldError: true,
		},
		"start after end": {
			sliceLength: 4,
			start:       3,
			end:         2,
			shouldError: true,
		},
	}

	for name, test := range tests {
		err := ValidateIndexes(test.sliceLength, test.start, test.end)
		isError := err != nil
		if isError != test.shouldError {
			if test.shouldError {
				t.Fatalf("Test \"%s\" did not error when it should have", name)
			} else {
				t.Fatalf("Test \"%s\" errored when it should not have", name)
			}
		}
	}
}

type stringSliceContainsTest struct {
	slice  []string
	value  string
	wanted bool
}

func TestStringSliceContains(t *testing.T) {
	tests := map[string]stringSliceContainsTest{
		"contains": {
			slice:  []string{"a", "b", "c"},
			value:  "a",
			wanted: true,
		},
		"does not contain": {
			slice:  []string{"a", "b", "c"},
			value:  "d",
			wanted: false,
		},
		"empty string": {
			slice:  []string{"a", "b", "c"},
			value:  "",
			wanted: false,
		},
		"repeated entries": {
			slice:  []string{"a", "a", "a"},
			value:  "a",
			wanted: true,
		},
	}

	for name, test := range tests {
		got := StringSliceContains(test.slice, test.value)
		if got != test.wanted {
			t.Fatalf(
				"test \"%s\" returned the wrong result. Got: %t, wanted: %t",
				name,
				got,
				test.wanted,
			)
		}
	}
}

type stringSliceEqualTest struct {
	slice1   []string
	slice2   []string
	expected bool
}

func TestStringSlicesEqual(t *testing.T) {
	tests := map[string]stringSliceEqualTest{
		"equal": {
			slice1:   []string{"hello", "world"},
			slice2:   []string{"hello", "world"},
			expected: true,
		},
		"both empty": {
			slice1:   []string{},
			slice2:   []string{},
			expected: true,
		},
		"one empty, one not": {
			slice1:   []string{"hello", "world"},
			slice2:   []string{},
			expected: false,
		},
		"different lengths": {
			slice1:   []string{"hello", "world"},
			slice2:   []string{"hello", "world", "foo"},
			expected: false,
		},
		"same length, different values": {
			slice1:   []string{"hello", "sir"},
			slice2:   []string{"hello", "world"},
			expected: false,
		},
	}

	for name, test := range tests {
		got := StringSlicesEqual(test.slice1, test.slice2)

		if got != test.expected {
			t.Fatalf(
				"Test \"%s\" returned wrong results. Wanted: %t, got: %t",
				name,
				test.expected,
				got,
			)
		}
	}
}

type appendStringIfNotInSliceTest struct {
	slice  []string
	s      string
	wanted []string
}

func TestAppendStringIfNotInSlice(t *testing.T) {
	tests := map[string]appendStringIfNotInSliceTest{
		"string already in slice": {
			slice:  []string{"hello"},
			s:      "hello",
			wanted: []string{"hello"},
		},
		"string not in slice": {
			slice:  []string{"hello"},
			s:      "world",
			wanted: []string{"hello", "world"},
		},
		"adding to empty slice": {
			slice:  []string{},
			s:      "hello",
			wanted: []string{"hello"},
		},
	}

	for name, test := range tests {
		got := AppendStringIfNotInSlice(test.slice, test.s)

		if !StringSlicesEqual(got, test.wanted) {
			t.Fatalf("Test \"%s\" failed. Wanted: %v, got: %v", name, test.wanted, got)
		}
	}
}
