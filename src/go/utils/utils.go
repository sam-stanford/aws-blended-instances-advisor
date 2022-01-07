package utils

import (
	"fmt"
	"os"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

const (
	WRITE_PERMISSION_CODE = 0755
)

// GenerateUuid generates a universally unique identifier (UUID),
// returned as a string.
func GenerateUuid() string {
	return uuid.NewString()
}

// StopProgramExecution ends the program's execution, printing the provided error
// and returning the provided error code.
func StopProgramExecution(err error, exitCode int) {
	fmt.Fprintf(os.Stderr, "Stopping program execution: %s\n", err.Error())
	os.Exit(exitCode)
}

// PrependToError formats and adds a string in front of a given error.
func PrependToError(err error, message string) error {
	return fmt.Errorf("%s: %s", message, err.Error())
}

// Checks whether start (inclusive) and end (exclusive) are valid indexes for
// the given lenght of a slice.
func ValidateIndexes(sliceLength, start, end int) error {
	if start >= end {
		return fmt.Errorf("provided indexes of %d and %d provide a subslice of 0 elements", start, end)
	}
	if start < 0 {
		return fmt.Errorf("provided start of %d is out of bounds for slice of length %d", start, sliceLength)
	}
	if end > sliceLength {
		return fmt.Errorf("provided end of %d is out of bounds for slice of length %d", end, sliceLength)
	}
	return nil
}

// CreateMockLogger creates a mock logger which can be
// passed as a zap.Logger during testing.
func CreateMockLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{}
	return cfg.Build()
}

// StringSliceContains returns true if the given slice contains the given
// value, and false otherewise.
func StringSliceContains(slice []string, value string) bool {
	for _, val := range slice {
		if val == value {
			return true
		}
	}
	return false
}

// StringSlicesEqual returns a Boolean with value representing whether the
// two provided slices are equal in terms of length and values.
func StringSlicesEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// AppendStringIfNotInSlice appends a string to a slice if the slice
// does not already contain the string.
func AppendStringIfNotInSlice(slice []string, s string) []string {
	inSlice := false
	for i := range slice {
		if slice[i] == s {
			inSlice = true
			break
		}
	}

	if !inSlice {
		return append(slice, s)
	}
	return slice
}
