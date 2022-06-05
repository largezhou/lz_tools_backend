package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"github.com/largezhou/lz_tools_backend/app/middleware"
)

func InitServer(ctx context.Context) *gin.Engine {
	if c.Debug {
		logger.Debug(ctx, "debug 模式运行")
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.MaxMultipartMemory = 10 << 20
	r.Use(middleware.SetRequestId(), middleware.Logger())

	return r
}
