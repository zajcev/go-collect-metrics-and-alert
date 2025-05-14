package middleware

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"sync"
	"time"
)

type data struct {
	http.ResponseWriter
	statusCode int
	contentLen int
}

var (
	logger     *zap.Logger
	loggerOnce sync.Once
)

func initLogger() {
	cfg := zap.NewProductionConfig()
	cfg.Sampling = nil
	cfg.OutputPaths = []string{"stdout"}

	cfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		MessageKey:    "msg",
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.EpochTimeEncoder,
		StacktraceKey: "",
	}

	var err error
	logger, err = cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
}

func ZapMiddleware(wrappedHandler http.Handler) http.Handler {
	loggerOnce.Do(initLogger)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		logger.Info("Request data",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
		)
		rw := newLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(rw, r)
		logger.Info("Response data",
			zap.Int("statusCode", rw.statusCode),
			zap.Int("contentLen", rw.contentLen),
			zap.Duration("duration", time.Since(now)),
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
