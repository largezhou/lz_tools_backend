package console

import (
	"github.com/largezhou/gin_starter/app/app_const"
	"github.com/urfave/cli/v2"
	"os"
)

// RunInCli 是否在运行 CLI 命令
func RunInCli() bool {
	args := os.Args
	return len(args) >= 2 && args[1] == app_const.CliKey
}

var App *cli.App

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
