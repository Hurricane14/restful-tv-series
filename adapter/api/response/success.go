package response

import (
	"encoding/json"
	"net/http"
)

type Success struct {
	statusCode int
	result     any
}

func NewSuccess(code int, result any) *Success {
	return &Success{
		statusCode: code,
		result:     result,
	}
}

func (s *Success) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s.statusCode)
	if s.result != nil {
		return json.NewEncoder(w).Encode(s.result)
	}
	return nil
}
