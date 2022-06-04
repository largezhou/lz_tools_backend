package helper

import (
	"github.com/largezhou/lz_tools_backend/app/config"
)

var shutdownFuncList []func()

// RegisterShutdownFunc 注册一个服务关闭时的回调函数
func RegisterShutdownFunc(f func()) {
	shutdownFuncList = append(shutdownFuncList, f)
}

// CallShutdownFunc 服务关闭时，执行所有回调函数
func CallShutdownFunc() {
	for _, f := range shutdownFuncList {
		f()
	}
}

// CheckAppKey 检查 app key
func CheckAppKey() {
	if len(config.Config.App.Key) < 32 {
		panic("APP_KEY 长度至少为 32 位")
	}
}
