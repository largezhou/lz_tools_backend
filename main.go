package main

import (
	"context"
	"github.com/largezhou/gin_starter/app/api"
	"github.com/largezhou/gin_starter/app/console"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/largezhou/gin_starter/app/config"
	"github.com/largezhou/gin_starter/app/helper"
	"github.com/largezhou/gin_starter/app/logger"
	"go.uber.org/zap"

	_ "github.com/largezhou/gin_starter/app/console"
	_ "github.com/largezhou/gin_starter/app/model"
)

var c = config.Config.App
var r *gin.Engine

func main() {
	if console.RunInCli() {
		return
	}

	defer helper.CallShutdownFunc()

	r = InitServer()
	api.InitRouter(r)

	srv := &http.Server{
		Addr:    c.Host + ":" + c.Port,
		Handler: r,
	}

	go func() {
		logger.Info("开始运行", zap.String("host", c.Host), zap.String("port", c.Port))
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
		logger.Error("服务关闭出错", zap.Error(err))
	}

	logger.Info("服务已关闭")
}
