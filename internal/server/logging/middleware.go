package logging

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

func NewMiddleware(wrappedHandler http.Handler) http.Handler {
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
		statusCode := rw.statusCode
		logger.Info(
			"Response data",
			zap.Int("statusCode", statusCode),
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
