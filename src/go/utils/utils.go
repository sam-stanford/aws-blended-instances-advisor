package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// TODO: Doc comments

func FileToString(filepath string) (string, error) {
	fileAsBytes, err := FileToBytes(filepath)
	return string(fileAsBytes), err
}

func FileToBytes(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

func GetCallerPath() (string, error) {
	return os.Getwd()
}

func CreateFilepath(pathComponents ...string) (string, error) {
	path := ""
	for idx, component := range pathComponents {
		path += component
		if idx != len(pathComponents)-1 {
			path += string(os.PathSeparator)
		}
	}
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
