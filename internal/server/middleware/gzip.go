package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/crypto"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"io"
	"log"
	"net/http"
	"strings"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

type responseWriter struct {
	http.ResponseWriter
	compressWriter *compressWriter
	wroteHeader    bool
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.wroteHeader {
		contentType := rw.Header().Get("Content-Type")
		if contentType == "application/json" || strings.Contains(contentType, "text/html") {
			rw.Header().Set("Content-Encoding", "gzip")
		}
		rw.wroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.compressWriter.Write(b)
}

type byteReadCloser struct {
	*bytes.Reader
}

func (b *byteReadCloser) Close() error {
	return nil
}

func newByteReadCloser(data []byte) io.ReadCloser {
	return &byteReadCloser{bytes.NewReader(data)}
}

func newCompressReaderFromBytes(data []byte) (*compressReader, error) {
	rc := newByteReadCloser(data)
	defer func(rc io.ReadCloser) {
		err := rc.Close()
		if err != nil {
		}
	}(rc)
	zr, err := gzip.NewReader(rc)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  rc,
		zr: zr,
	}, nil
}
func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Request
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body from request : %v", err)
		}
		if sendsGzip {
			if config.GetCryptoKey() != "" {
				key, errKey := crypto.LoadPrivateKey(config.GetCryptoKey())
				if errKey != nil {
					log.Printf("Error load public key : %v", errKey)
				}
				decrypted, errDecrypt := crypto.Decrypt(key, body)
				if errDecrypt != nil {
					log.Printf("Error decrypt data : %v", errDecrypt)
				}
				body = decrypted
			}
			cr, errCompress := newCompressReaderFromBytes(body)
			if errCompress != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer func() {
				if err = cr.Close(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}()
		}

		// Response
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			defer func() {
				if err := cw.Close(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}()
			rw := &responseWriter{ResponseWriter: w, compressWriter: cw}
			h.ServeHTTP(rw, r)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
