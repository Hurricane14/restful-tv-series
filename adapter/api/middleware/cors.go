package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func CORS(next http.Handler) http.Handler {
	return cors.Default().Handler(next)
}
