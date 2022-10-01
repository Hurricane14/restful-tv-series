package middleware

import (
	"bytes"
	"io"
	"net/http"
	"series/adapter/logger"
	"time"
)

type loggedResponseWriter struct {
	code int
	w    http.ResponseWriter
	body *bytes.Buffer
	mw   io.Writer
}

func (l *loggedResponseWriter) Write(b []byte) (int, error) {
	return l.mw.Write(b)
}

func (l *loggedResponseWriter) WriteHeader(code int) {
	l.code = code
	l.w.WriteHeader(code)
}

func (l *loggedResponseWriter) Header() http.Header {
	return l.w.Header()
}

func newLoggedResponseWriter(w http.ResponseWriter) *loggedResponseWriter {
	buff := &bytes.Buffer{}
	return &loggedResponseWriter{
		w:    w,
		body: buff,
		mw:   io.MultiWriter(w, buff),
	}
}

func Logging(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log.Infof("%s Method: %s Path: %s", start.Format(time.RFC822), r.Method, r.URL)

			lw := newLoggedResponseWriter(w)
			next.ServeHTTP(lw, r)

			end := time.Now()

			var logFn func(string, ...any)
			switch {
			case lw.code >= http.StatusInternalServerError:
				logFn = log.Errorf
			case lw.code >= http.StatusBadRequest:
				logFn = log.Warnf
			default:
				logFn = log.Infof
			}

			logFn(
				"%s Method: %s Path: %s Status: %s Elapsed time: %s",
				end.Format(time.RFC822),
				r.Method,
				r.URL,
				http.StatusText(lw.code),
				end.Sub(start),
			)
		})
	}
}
