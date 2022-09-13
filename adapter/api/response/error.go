package response

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	statusCode int
	errors     []string
}

func NewError(code int, errors ...string) *Error {
	return &Error{
		statusCode: code,
		errors:     errors,
	}
}

func (e *Error) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.statusCode)
	if len(e.errors) == 0 {
		return nil
	}
	return json.NewEncoder(w).Encode(
		map[string][]string{
			"errors": e.errors,
		},
	)
}
