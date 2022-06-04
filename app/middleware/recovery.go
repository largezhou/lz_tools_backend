package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"go.uber.org/zap"
)

// Recovery panic 处理
func Recovery(formatter func(*gin.Context, any)) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("异常", zap.Any("app_error", err))
				formatter(c, err)
			}
		}()
		c.Next()
	}
}
