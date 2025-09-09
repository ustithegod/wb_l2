package event

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrInternalServerError = errors.New("internal server error")
	ErrServiceUnavailable  = errors.New("service is temporary unavailable")
)

func GetStatusCode(err error) int {
	if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest
	} else if errors.Is(err, ErrServiceUnavailable) {
		return http.StatusServiceUnavailable
	}

	return http.StatusInternalServerError
}
