package middleware

import (
	"errors"
	"net/http"
)

func Root2Index(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "/index" {
		return errors.New("Please access the root directory `/`")
	}
	if r.URL.Path == "/" {
		r.URL.Path = "/index"
	}
	return nil
}
