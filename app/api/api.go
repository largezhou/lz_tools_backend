package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/largezhou/lz_tools_backend/app/app_const"
	"github.com/largezhou/lz_tools_backend/app/app_error"
	"github.com/largezhou/lz_tools_backend/app/config"
	"github.com/largezhou/lz_tools_backend/app/middleware"
	"github.com/largezhou/lz_tools_backend/app/model/user_model"
	"net/http"
	"runtime/debug"
)

func InitRouter(r *gin.Engine) {
	codeController := NewCodeController()

	{
		g := getApiGroup(r).Use(
			middleware.Cors(getApiGroup(r)),
			middleware.Recovery(failAny),
			UsernameAuth(),
		)

		g.POST("/get-code-list", codeController.GetCodeList)
		g.POST("/save-code", codeController.SaveCode)
		g.POST("/delete-code", codeController.DeleteCode)
	}
}

func getApiGroup(r *gin.Engine) *gin.RouterGroup {
	return r.Group("/api")
}

func failAny(ctx *gin.Context, err any) {
	realErr, ok := err.(error)
	if ok {
		fail(ctx, realErr)
	} else {
		handleDefaultError(ctx, err)
	}
}

func fail(ctx *gin.Context, err error) {
	switch {
	case errors.As(err, &validator.ValidationErrors{}):
		response(ctx, app_error.InvalidParameter, err.Error(), nil, nil)
	case errors.As(err, &app_error.Error{}):
		e := err.(app_error.Error)
		response(ctx, e.Code, e.Msg, nil, nil)
	default:
		handleDefaultError(ctx, err)
	}
}

func failWith(ctx *gin.Context, code int, msg string) {
	response(ctx, code, msg, nil, nil)
}

func handleDefaultError(ctx *gin.Context, err any) {
	var msg string
	fields := gin.H{}

	if config.Config.App.Debug {
		msg = fmt.Sprintf("%v", err)
		fields["trace"] = string(debug.Stack())
	} else {
		msg = http.StatusText(http.StatusInternalServerError)
	}

	response(ctx, app_error.UnknownErr, msg, nil, fields)
}

func ok(ctx *gin.Context, data any, msg string) {
	response(ctx, app_error.StatusOk, msg, data, nil)
}

func response(ctx *gin.Context, code int, msg string, data any, fields gin.H) {
	resp := gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	}

	for k, v := range fields {
		resp[k] = v
	}

	ctx.JSON(http.StatusOK, resp)
}

// getAuthUser 获取已登录用户
func getAuthUser(ctx *gin.Context) (*user_model.User, error) {
	var user *user_model.User
	userAny, ok := ctx.Get(app_const.AuthUserKey)
	if !ok {
		return user, fmt.Errorf("没有用户信息")
	}

	user, ok = userAny.(*user_model.User)
	if !ok {
		return user, fmt.Errorf("无法获取用户信息")
	}

	return user, nil
}
