package utils

import (
	"io"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

// TODO: Doc & test

func WriteHttpErrorResponse(
	w http.ResponseWriter,
	requestId string,
	err error,
	errCode int,
	logger *zap.Logger,
) {
	http.Error(w, err.Error(), errCode)

	logger.Error(
		"responded to request with error",
		zap.String("reqId", requestId),
		zap.Int("responseCode", errCode),
		zap.Error(err),
	)
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
