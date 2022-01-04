package utils

import (
	"fmt"
	"os"
	"reflect"

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

// AnyFieldsAreEmpty checks if any fields of a given interface are empty.
//
// If at least one field is considered empty, the first empty field is
// also returned.
func AnyFieldsAreEmpty(i interface{}) (bool, string) {
	elem := reflect.ValueOf(i).Elem()
	return anyFieldsAreEmptyHelper(elem)
}

func anyFieldsAreEmptyHelper(elem reflect.Value) (bool, string) {
	for fieldIndex := 0; fieldIndex < elem.NumField(); fieldIndex += 1 {

		field := elem.Field(fieldIndex)

		if field.Kind() == reflect.Struct {
			empty, field := anyFieldsAreEmptyHelper(field)
			if empty {
				return true, field
			}
		} else if !isNumericOrBoolType(field.Kind()) && field.IsZero() {
			return true, elem.Type().Field(fieldIndex).Name
		}
	}

	return false, ""
}

func isNumericOrBoolType(k reflect.Kind) bool {
	types := []reflect.Kind{
		reflect.Float32,
		reflect.Float64,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Bool,
	}

	for _, t := range types {
		if k == t {
			return true
		}
	}
	return false
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
