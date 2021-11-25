package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"go.uber.org/zap"
)

const (
	WRITE_PERMISSION_CODE = 0755
)

// TODO: Split into multiple files

// TODO: Doc comments

func FileToString(filepath string) (string, error) {
	fileAsBytes, err := FileToBytes(filepath)
	return string(fileAsBytes), err
}

func FileToBytes(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

func WriteBytesToFile(data []byte, filepath string) error {
	return os.WriteFile(filepath, data, WRITE_PERMISSION_CODE)
}

// TODO: Test
func WriteStringToFile(data string, filepath string) error {
	return WriteBytesToFile([]byte(data), filepath)
}

// TODO: Test
func DeleteFile(filepath string) error {
	return os.Remove(filepath)
}

// TODO: Test
func FileExists(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func GetCallerPath() (string, error) {
	return os.Getwd()
}

// TODO: Rename
func CreateFilepath(pathComponents ...string) (string, error) {
	path := ""
	for idx, component := range pathComponents {
		path += component
		if idx != len(pathComponents)-1 {
			path += string(os.PathSeparator)
		}
	}
	return AbsoluteFilepath(path)
}

// TODO: Test
func AbsoluteFilepath(path string) (string, error) {
	return filepath.Abs(path)
}

func DownloadFile(url string, filepath string) error {
	// TODO: Check if file has changed using HEAD & checking for Last-Modified field

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func GetHttpLastModified(url string) (time.Time, error) {
	header, err := GetHttpHeader(url)
	if err != nil {
		return time.Now(), err
	}

	lastModifiedString := header.Get("Last-Modified")
	return time.Parse(time.RFC1123, lastModifiedString)
}

func GetHttpHeader(url string) (http.Header, error) {
	head, err := http.Head(url)
	if err != nil {
		return nil, err
	}
	return head.Header, nil
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
