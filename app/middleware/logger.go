package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/largezhou/lz_tools_backend/app/app_const"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"go.uber.org/zap"
	"time"
)

// Logger 记录请求，添加 requestId
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.NewString()
		c.Set(app_const.RequestIdKey, requestId)
		logger.Logger = logger.Logger.With(zap.String("requestId", requestId))

		start := time.Now()

		logger.Info(
			"request",
			zap.String("clientIp", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
		)

		c.Next()

		logger.Info(
			"response",
			zap.Duration("cost", time.Now().Sub(start)),
			zap.Int("code", c.Writer.Status()),
		)
	}
}
