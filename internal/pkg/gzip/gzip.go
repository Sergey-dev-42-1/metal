package gzip

import (
	"compress/gzip"

	"fmt"
	"io"

	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	gz *gzip.Writer
}

func (zw *gzipWriter) Write(p []byte) (int, error) {
	return zw.gz.Write(p)
}

func (zw *gzipWriter) WriteString(s string) (int, error) {
	return zw.gz.Write([]byte(s))
}
func (zw *gzipWriter) WriteHeader(statusCode int) {
	zw.Header().Set("Content-Encoding", "gzip")
	zw.ResponseWriter.WriteHeader(statusCode)
}

type gzipReader struct {
	r  io.ReadCloser
	gr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*gzipReader, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r:  r,
		gr: gr,
	}, nil
}

func (zr *gzipReader) Read(p []byte) (int, error) {
	return zr.gr.Read(p)
}

func (zr *gzipReader) Close() error {
	if err := zr.gr.Close(); err != nil {
		return err
	}
	return zr.r.Close()
}

func GzipHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		receivedGzip := strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip")
		if receivedGzip {
			gz, err := newCompressReader(c.Request.Body)
			if err != nil {
				fmt.Println(err)
				_ = c.AbortWithError(400, err)
				return
			}
			// c.Request.Header.Del("Content-Encoding")
			// c.Request.Header.Del("Content-Length")
			c.Request.Body = gz
			defer gz.Close()
		}

		acceptsGzip := strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip")
		if acceptsGzip {
			gz := gzip.NewWriter(c.Writer)
			w := &gzipWriter{c.Writer, gz}
			c.Writer = w
			defer func() {
				gz.Close()
			}()

		}
		c.Next()
	}
}

// type gzipReader struct {
// 	gin.
// 	body *bytes.Buffer
// }

// func (w bodyLogWriter) Write(b []byte) (int, error) {
// 	w.body.Write(b)
// 	return w.ResponseWriter.Write(b)
// }
