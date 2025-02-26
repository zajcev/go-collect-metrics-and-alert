package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type data struct {
	http.ResponseWriter
	statusCode int
	contentLen int
}

func ZapMiddleware(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		logger, err := zap.NewProduction()
		defer logger.Sync()
		if err != nil {
			logger.Fatal(err.Error())
			return
		}
		logger.Info("Request data",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", time.Since(now)))
		rw := newLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(rw, r)
		logger.Info(
			"Response data",
			zap.Int("statusCode", rw.statusCode),
			zap.Int("contentLen", rw.contentLen),
		)
	})
}

func newLoggingResponseWriter(w http.ResponseWriter) *data {
	return &data{w, http.StatusOK, -1}
}

func (d *data) WriteHeader(code int) {
	d.statusCode = code
	d.ResponseWriter.WriteHeader(code)
}

func (d *data) Write(b []byte) (int, error) {
	d.contentLen += len(b)
	return d.ResponseWriter.Write(b)
}
