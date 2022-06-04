package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 记录请求，添加 requestId
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		rawData, _ := ctx.GetRawData()
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))

		logger.Info(
			ctx,
			"request",
			zap.String("clientIp", ctx.ClientIP()),
			zap.String("path", ctx.Request.URL.Path),
			zap.String("query", ctx.Request.URL.RawQuery),
			zap.ByteString("data", rawData),
		)

		blw := &bodyLogWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}
		ctx.Writer = blw

		ctx.Next()

		logger.Info(
			ctx,
			"response",
			zap.Duration("cost", time.Now().Sub(start)),
			zap.Int("code", ctx.Writer.Status()),
			zap.ByteString("data", blw.body.Bytes()),
		)
	}
}
