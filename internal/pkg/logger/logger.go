package logger

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger = zap.NewNop().Sugar()

func New(level string) error {

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	fmt.Println(lvl)
	cfg := zap.NewDevelopmentConfig()

	zl, err := cfg.Build()
	defer zl.Sync()
	if err != nil {
		return err
	}
	log = zl.Sugar()
	return nil
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uri := c.Request.RequestURI
		method := c.Request.Method
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		duration := time.Since(start)
		log.Infoln(
			"uri", uri,
			"duration", duration,
			"method", method,
			"status", c.Writer.Status(),
			"size", c.Writer.Size(),
		)
	}
}
