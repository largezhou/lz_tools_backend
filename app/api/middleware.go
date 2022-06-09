package api

import (
	"github.com/gin-gonic/gin"
	"github.com/largezhou/lz_tools_backend/app/app_const"
	"github.com/largezhou/lz_tools_backend/app/app_error"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"github.com/largezhou/lz_tools_backend/app/model/user_model"
	"go.uber.org/zap"
)

// UsernameAuth 只需通过请求头中的用户名即可认证通过
func UsernameAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetHeader("Username")
		if username == "" {
			response(ctx, app_error.AuthFail, "没有用户信息", nil, nil)
			ctx.Abort()
			return
		}

		if len(username) < 3 || len(username) > 20 {
			response(ctx, app_error.InvalidParameter, "用户名长度为 3-20", nil, nil)
			ctx.Abort()
			return
		}

		user := &user_model.User{Username: username}
		var err error
		if user, err = user_model.UpdateOrCreateUserByUserInfo(ctx, user); err != nil {
			logger.Warn(ctx, "创建或更新用户失败", zap.Error(err))
			fail(ctx, err)
			ctx.Abort()
			return
		}

		ctx.Set(app_const.AuthUserKey, user)
		ctx.Next()
	}
}
