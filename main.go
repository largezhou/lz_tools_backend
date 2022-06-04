package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/largezhou/lz_tools_backend/app/api"
	"github.com/largezhou/lz_tools_backend/app/app_const"
	"github.com/largezhou/lz_tools_backend/app/console"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/largezhou/lz_tools_backend/app/config"
	"github.com/largezhou/lz_tools_backend/app/helper"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"go.uber.org/zap"

	_ "github.com/largezhou/lz_tools_backend/app/console"
)

var c = config.Config.App
var r *gin.Engine

func main() {
	ctx := context.WithValue(context.Background(), app_const.RequestIdKey, uuid.NewString())

	if console.RunInCli() {
		return
	}

	defer helper.CallShutdownFunc(ctx)

	r = InitServer(ctx)
	api.InitRouter(r)

	srv := &http.Server{
		Addr:    c.Host + ":" + c.Port,
		Handler: r,
	}

	go func() {
		logger.Info(ctx, "开始运行", zap.String("host", c.Host), zap.String("port", c.Port))
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(ctx, "服务关闭出错", zap.Error(err))
	}

	logger.Info(ctx, "服务已关闭")
}
