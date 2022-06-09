package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/largezhou/lz_tools_backend/app/api"
	"github.com/largezhou/lz_tools_backend/app/app_const"
	"github.com/largezhou/lz_tools_backend/app/core"
)

func main() {
	ctx := context.WithValue(context.Background(), app_const.RequestIdKey, uuid.NewString())

	app := core.Get(ctx)

	if !core.RunInConsole(ctx) {
		api.InitRouter(app.Engine)
	}

	app.Run(ctx)
}
