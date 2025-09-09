package http

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]any

func writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	response, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	response = append(response, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)

	return nil
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	writeJson(w, status, envelope{"error": message}, nil)
}
