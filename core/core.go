package core

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/largezhou/lz_tools_backend/app/app_const"
	"github.com/largezhou/lz_tools_backend/app/config"
	"github.com/largezhou/lz_tools_backend/app/console"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"github.com/largezhou/lz_tools_backend/app/middleware"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	*gin.Engine
	shutdownFuncList []ShutdownFunc
	console          *cli.App
	args             []string
	rawArgs          []string
}

type ShutdownFunc func(context.Context, *App)

var app *App

// Get 获取全局唯一 App 实例
func Get(ctx context.Context) *App {
	if app == nil {
		app = New(ctx)
	}

	return app
}

// New 总是新建一个 App 实例
func New(ctx context.Context) *App {
	return &App{
		Engine:           initServer(ctx),
		shutdownFuncList: []ShutdownFunc{},
		console: &cli.App{
			Commands: console.Commands,
		},
		rawArgs: os.Args,
		// 去掉了 cli 参数之后的 args
		args: initArgs(ctx),
	}
}

func RunInConsole(ctx context.Context) bool {
	args := os.Args
	return len(args) >= 2 && args[1] == app_const.CliKey
}

func initServer(ctx context.Context) *gin.Engine {
	if RunInConsole(ctx) {
		return nil
	}

	c := config.Config.App
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

func initArgs(ctx context.Context) []string {
	var args []string
	// 第二个参数为 CLI 入口，复制并删除
	for i, arg := range os.Args {
		if i != 1 {
			args = append(args, arg)
		}
	}

	return args
}

// RegisterShutdownFunc 注册一个服务关闭时的回调函数
func (a *App) RegisterShutdownFunc(f ShutdownFunc) {
	a.shutdownFuncList = append(a.shutdownFuncList, f)
}

// callShutdownFunc 服务关闭时，执行所有回调函数
func (a *App) callShutdownFunc(ctx context.Context) {
	for _, f := range a.shutdownFuncList {
		f(ctx, a)
	}
}

func (a *App) Run(ctx context.Context) {
	defer a.callShutdownFunc(ctx)

	if RunInConsole(ctx) {
		a.runConsole(ctx)
	} else {
		a.runServer(ctx)
	}
}

func (a *App) runServer(ctx context.Context) {
	c := config.Config.App
	srv := &http.Server{
		Addr:    c.Host + ":" + c.Port,
		Handler: a.Engine,
	}

	go func() {
		logger.Info(ctx, "开始运行", zap.String("host", c.Host), zap.String("port", c.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(ctx, "服务关闭出错", zap.Error(err))
	}

	logger.Info(ctx, "服务已关闭")
}

func (a *App) runConsole(ctx context.Context) {
	if err := a.console.Run(a.args); err != nil {
		logger.Error(ctx, "命令行执行失败", zap.Error(err))
		panic(err)
	}
}
