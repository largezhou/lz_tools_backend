package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/largezhou/lz_tools_backend/app/config"
	"github.com/largezhou/lz_tools_backend/app/middleware"
	"net/http"
	"runtime/debug"
)

const (
	statusOk         = 0
	unknownErr       = 50000
	badRequest       = 40000
	authFail         = 40001
	invalidParameter = 40022
)

func InitRouter(r *gin.Engine) {
	commonMiddlewares := []gin.HandlerFunc{
		middleware.Cors(getApiGroup(r)),
		middleware.Recovery(errorToJsonResponse),
	}

	{
		g := getApiGroup(r).Use(commonMiddlewares...)

		g.POST("/login", getJwtMiddleware().LoginHandler)
		g.POST("/get-wechat-auth-url", getWechatAuthUrl)
	}

	{
		g := getApiGroup(r).Use(commonMiddlewares...).Use(getJwtMiddleware().MiddlewareFunc())

		g.POST("/get-code", getCode)
		g.POST("/get-code-list", getCodeList)
	}
}

func getApiGroup(r *gin.Engine) *gin.RouterGroup {
	return r.Group("/api")
}

// errorToJsonResponse 把 panic 处理成 json 数据
func errorToJsonResponse(ctx *gin.Context, err any) {
	var msg string
	fields := gin.H{}

	if config.Config.App.Debug {
		msg = fmt.Sprintf("%v", err)
		fields["trace"] = string(debug.Stack())
	} else {
		msg = http.StatusText(http.StatusInternalServerError)
	}

	response(ctx, unknownErr, msg, nil, fields)
}

func ok(ctx *gin.Context, data any, msg string) {
	response(ctx, statusOk, msg, data, nil)
}

func fail(ctx *gin.Context, code int, msg string) {
	response(ctx, code, msg, nil, nil)
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
