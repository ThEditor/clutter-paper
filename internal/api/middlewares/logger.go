package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(
			"[INC] %s \"%s %s %s\" from %s\n",
			time.Now().Format("2006/01/02 15:04:05"),
			r.Method,
			r.URL.String(),
			r.Proto,
			r.RemoteAddr,
		)

		start := time.Now()

		ww := newResponseWriterWrapper(w)

		next.ServeHTTP(ww, r)

		fmt.Printf(
			"[OUT] %s \"%s %s %s\" from %s - %d %dB in %v\n",
			time.Now().Format("2006/01/02 15:04:05"),
			r.Method,
			r.URL.String(),
			r.Proto,
			r.RemoteAddr,
			ww.statusCode,
			ww.bytesWritten,
			time.Since(start),
		)
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{w, http.StatusOK, 0}
}

func (ww *responseWriterWrapper) WriteHeader(code int) {
	ww.statusCode = code
	ww.ResponseWriter.WriteHeader(code)
}

func (ww *responseWriterWrapper) Write(b []byte) (int, error) {
	n, err := ww.ResponseWriter.Write(b)
	ww.bytesWritten += n
	return n, err
}
