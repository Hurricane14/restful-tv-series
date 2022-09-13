package response

import "net/http"

type Response interface {
	Send(w http.ResponseWriter) error
}
