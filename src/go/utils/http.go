package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
)

// DownloadFile makes a GET request to the given URL, saving the
// fetched data to the given filepath.
func DownloadFile(url string, filepath string) error {
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

// AddJsonContentTypeHeader adds the appropriate header to a request to set the
// response content type to JSON.
func AddJsonContentTypeHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// AddCorsHeader adds the appropriate header(s) to a response to allow cross-origin access
// for the given domains if they are allowed, returning an error otherwise.
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
