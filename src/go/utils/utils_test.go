package utils

import (
	"errors"
	"os"
	"testing"
)

const (
	TEST_TXT_FILEPATH     = "../../../assets/test/test.txt"
	TEST_TXT_FILE_CONTENT = "TEST TEXT"
	TEST_TEMP_FILEPATH    = "../../../assets/test/temp.txt"
)

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

func TestWriteBytesToFileForNewFile(t *testing.T) {
	_, err := os.Stat(TEST_TEMP_FILEPATH)
	if err == nil {
		t.Fatalf("Cannot test for new file as file already exists. Filepath: %s", TEST_TEMP_FILEPATH)
	}

	b := []byte("TESTING123")
	err = WriteBytesToFile(b, TEST_TEMP_FILEPATH)
	if err != nil {
		t.Fatalf("Failed to write bytes to file: %s", err.Error())
	}

	err = DeleteFile(TEST_TEMP_FILEPATH)
	if err != nil {
		t.Fatalf(
			"Failed to delete file. Manual deletion likely required. Filepath: %s, Error: %s",
			TEST_TEMP_FILEPATH,
			err.Error(),
		)
	}
}

func TestWriteBytesToFileForExistingFile(t *testing.T) {
	file, err := os.Create(TEST_TEMP_FILEPATH)
	if err != nil {
		t.Fatalf("Failed to create file in test setup: %s", err.Error())
	}
	err = file.Close()
	if err != nil {
		t.Fatalf("Failed to close file in test setup: %s", err.Error())
	}

	_, err = os.Stat(TEST_TEMP_FILEPATH)
	if err != nil {
		t.Fatalf(
			"Cannot test for existing file as file does not exist. Filepath: %s, Error: %s",
			TEST_TEMP_FILEPATH,
			err.Error(),
		)
	}

	b := []byte("TESTING123")
	err = WriteBytesToFile(b, TEST_TEMP_FILEPATH)
	if err != nil {
		t.Fatalf("Failed to write bytes to file: %s", err.Error())
	}

	err = DeleteFile(TEST_TEMP_FILEPATH)
	if err != nil {
		t.Fatalf(
			"Failed to delete file. Manual deletion likely required. Filepath: %s, Error: %s",
			TEST_TEMP_FILEPATH,
			err.Error(),
		)
	}
}

func TestCreateFilepath(t *testing.T) {
	// TODO
}

func TestDownloadFile(t *testing.T) {
	// TODO
}

func TestFileExists(t *testing.T) {
	exists, err := FileExists(TEST_TXT_FILEPATH)
	if err != nil {
		t.Fatalf("Error when checking if file exists: %s", err.Error())
	}
	if !exists {
		t.Fatalf("FileExists returned false for existing file. Wanted: %t, got: %t", true, exists)
	}

	exists, err = FileExists("../FILE_DOES_NOT_EXIST.txt")
	if err != nil {
		t.Fatalf("Error when checking if file exists: %s", err.Error())
	}
	if exists {
		t.Fatalf("FileExists returned true for non-existent file. Wanted: %t, got: %t", true, exists)
	}
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

func TestCreateMockLogger(t *testing.T) {
	logger, err := CreateMockLogger()
	if err != nil {
		t.Fatalf("Error returned when creating logger: %s", err.Error())
	}
	logger.Info("") // Should not throw error
}

// TODO: Test empty field checker
