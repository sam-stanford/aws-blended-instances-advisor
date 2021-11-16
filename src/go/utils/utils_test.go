package utils

import (
	"errors"
	"testing"
)

const (
	TEST_TXT_FILEPATH     = "../../../assets/test/test.txt"
	TEST_TXT_FILE_CONTENT = "TEST TEXT"
)

type testStruct struct {
	value int
}

func TestFileToBytes(t *testing.T) {
	bytes, err := FileToBytes(TEST_TXT_FILEPATH)
	if err != nil {
		t.Fatalf("Error thrown when reading file: %s", err.Error())
	}
	if string(bytes) != TEST_TXT_FILE_CONTENT {
		t.Fatalf("Read bytes do not match file content. Wanted: %s, got: %s", TEST_TXT_FILE_CONTENT, string(bytes))
	}
}

func TestFileToBytesThrowsErrorForInvalidFilepath(t *testing.T) {
	_, err := FileToBytes("NOT A FILEPATH")
	if err == nil {
		t.Fatalf("Error was not thrown for an invalid filepath")
	}
}

func TestFileToString(t *testing.T) {
	content, err := FileToString(TEST_TXT_FILEPATH)
	if err != nil {
		t.Fatalf("Error thrown when reading file: %s", err.Error())
	}
	if content != TEST_TXT_FILE_CONTENT {
		t.Fatalf("Read string does not match file content. Wanted: %s, got: %s", TEST_TXT_FILE_CONTENT, content)
	}
}

func TestFileToStringThrowsErrorForInvalidFilepath(t *testing.T) {
	_, err := FileToString("NOT A FILEPATH")
	if err == nil {
		t.Fatalf("Error was not thrown for an invalid filepath")
	}
}

func TestCreateFilepath(t *testing.T) {
	// TODO
}

func TestDownloadFile(t *testing.T) {
	// TODO
}

func TestGetHttpLastModified(t *testing.T) {
	// TODO
}

func TestGetHttpHeader(t *testing.T) {
	// TODO
}

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
