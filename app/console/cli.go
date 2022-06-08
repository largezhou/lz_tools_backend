package console

import (
	"context"
	"github.com/google/uuid"
	"github.com/largezhou/lz_tools_backend/app/app_const"
	"github.com/largezhou/lz_tools_backend/app/model"
	"github.com/urfave/cli/v2"
	"os"
)

// RunInCli 是否在运行 CLI 命令
func RunInCli() bool {
	args := os.Args
	return len(args) >= 2 && args[1] == app_const.CliKey
}

var App *cli.App
var ctx = context.WithValue(context.Background(), app_const.RequestIdKey, uuid.NewString())
var db = model.DB.WithContext(ctx)

func init() {
	if !RunInCli() {
		return
	}

	App = &cli.App{
		Commands: commands,
	}

	var args []string
	// 第二个参数为 CLI 入口，复制并删除
	for i, arg := range os.Args {
		if i != 1 {
			args = append(args, arg)
		}
	}

	if err := App.Run(args); err != nil {
		panic(err)
	}
}
