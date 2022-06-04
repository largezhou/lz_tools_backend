package helper

import (
	"context"
	"github.com/largezhou/lz_tools_backend/app/config"
)

type ShutdownFunc func(ctx context.Context)

var shutdownFuncList []ShutdownFunc

// RegisterShutdownFunc 注册一个服务关闭时的回调函数
func RegisterShutdownFunc(f ShutdownFunc) {
	shutdownFuncList = append(shutdownFuncList, f)
}

// CallShutdownFunc 服务关闭时，执行所有回调函数
func CallShutdownFunc(ctx context.Context) {
	for _, f := range shutdownFuncList {
		f(ctx)
	}
}

// CheckAppKey 检查 app key
func CheckAppKey() {
	if len(config.Config.App.Key) < 32 {
		panic("APP_KEY 长度至少为 32 位")
	}
}
