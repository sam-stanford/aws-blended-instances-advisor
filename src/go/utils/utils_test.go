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

// TODO: Test empty field checker
