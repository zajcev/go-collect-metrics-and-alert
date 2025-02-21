package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &compressReader{
		ReadCloser: r,
		zr:         zr,
	}, nil
}

func (c *compressReader) Read(p []byte) (int, error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.zr.Close(); err != nil {
		return err
	}
	return c.ReadCloser.Close()
}

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			wrw := &responseWriterWrapper{ResponseWriter: w}
			h.ServeHTTP(wrw, r)

			contentType := wrw.Header().Get("Content-Type")
			if contentType == "application/json" || contentType == "text/html" {
				w.Header().Set("Content-Encoding", "gzip")
				gz := gzip.NewWriter(w)
				defer gz.Close()
				gz.Write([]byte(wrw.body.String()))
			} else {
				for key, values := range wrw.Header() {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}
				w.Write([]byte(wrw.body.String()))
			}
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	body       strings.Builder
	statusCode int
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}
