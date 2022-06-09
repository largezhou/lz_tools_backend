package console

import (
	"context"
	"github.com/google/uuid"
	"github.com/largezhou/lz_tools_backend/app/app_const"
	"github.com/largezhou/lz_tools_backend/app/model"
)

var ctx = context.WithValue(context.Background(), app_const.RequestIdKey, uuid.NewString())
var db = model.DB.WithContext(ctx)
