package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

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
