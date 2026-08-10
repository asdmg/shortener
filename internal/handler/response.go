package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return
	}
}

func writeError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {
	writeJSON(
		w,
		status,
		ErrorResponse{
			Error:   code,
			Message: message,
		},
	)
}
