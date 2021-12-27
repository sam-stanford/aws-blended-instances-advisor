package utils

import (
	"fmt"
	"math"
	"os"
	"reflect"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

const (
	WRITE_PERMISSION_CODE = 0755
)

// TODO: Split into multiple files

// TODO: Doc comments

func GenerateUuid() string {
	return uuid.NewString()
}

func StopProgramExecution(err error, exitCode int) {
	fmt.Fprintf(os.Stderr, "Stopping program execution: %s\n", err.Error())
	os.Exit(exitCode)
}

func PrependToError(err error, message string) error {
	return fmt.Errorf("%s: %s", message, err.Error())
}

// Checks whether start (inclusive) and end (exclusive) are valid indexes for
// the given lenght of a slice.
// TODO: Test
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

func CreateMockLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{}
	return cfg.Build()
}

// TODO: Doc
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

// TODO: Test & Doc
func FloatsEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.001 // TODO: Make epsilon const
}

func MinOfInts(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxOfInts(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinOfFloats(a, b float64) float64 {
	return math.Min(a, b)
}

func MaxOfFloats(a, b float64) float64 {
	return math.Max(a, b)
}

func StringSliceContains(slice []string, value string) bool {
	for _, val := range slice {
		if val == value {
			return true
		}
	}
	return false
}
