package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

// TODO: Doc & test

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

// Adds the appropriate header to set the response content type to JSON
func AddJsonContentTypeHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// Adds the appropriate header(s) to allow cross-origin access for the given
// domains if they are allowed, returning an error otherwise.
func AddCorsHeader(w http.ResponseWriter, r *http.Request, allowedDomains []string) error {
	if r.Header == nil || r.Header["Origin"] == nil || len(r.Header["Origin"]) == 0 {
		return errors.New("origin header not provided on request")
	}
	origin := r.Header["Origin"][0]

	if !StringSliceContains(allowedDomains, origin) {
		return errors.New("origin is not allowed")
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	return nil
}
