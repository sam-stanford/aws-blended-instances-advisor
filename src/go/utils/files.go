package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileToString reads the content of the file at the given
// filepath into a string, returning an error if one is
// encountered while reading the file.
func FileToString(filepath string) (string, error) {
	fileAsBytes, err := FileToBytes(filepath)
	return string(fileAsBytes), err
}

// FileToBytes reads the content of a file at the given
// filepath into a byte array, returning an error if one
// is encountered while reading the file.
func FileToBytes(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

// WriteBytesToFile writes a byte array to the given
// filepath, returning an error if one is encountered while
// writing the to the file.
func WriteBytesToFile(data []byte, filepath string) error {
	return os.WriteFile(filepath, data, WRITE_PERMISSION_CODE)
}

// WriteStringToFile writes a string to the given filepath,
// returning an error if one is encountered while writing to
// the file.
func WriteStringToFile(data string, filepath string) error {
	return WriteBytesToFile([]byte(data), filepath)
}

// DeleteFile removes the file at the given filepath, returning
// an error if one is enocuntered while interacting with the filesystem.
func DeleteFile(filepath string) error {
	return os.Remove(filepath)
}

// FileExists checks whether a file exists at the given filepath and
// returns a boolean to convey the results.
//
// An error is returned if one is encountered while interacting with
// the filesystem.
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

// CreateFilepath forms a string filepath from a given set of components,
// adding system-specific path separators where appropriate.
//
// An error is returned if the generated filepath is invalid.
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

// AbsoluteFilepath returns an absolute representation of the given
// filepath, returning an error if the filepath is invalid.
func AbsoluteFilepath(path string) (string, error) {
	return filepath.Abs(path)
}
