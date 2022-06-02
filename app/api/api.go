package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/largezhou/gin_starter/app/app_const"
	"github.com/largezhou/gin_starter/app/config"
	"github.com/largezhou/gin_starter/app/middleware"
	"net/http"
	"runtime/debug"
)

const (
	statusOk   = 0
	unknownErr = 50000
	authFail   = 40001
)

func InitRouter(r *gin.Engine) {
	g := r.Group("/api").
		Use(
			middleware.Recovery(errorToJsonResponse),
			middleware.ApiAuth(),
		)

	g.POST("/get-code", GetCode)
	g.POST("/get-code-list", GetCodeList)
}

// errorToJsonResponse 把 panic 处理成 json 数据
func errorToJsonResponse(c *gin.Context, err any) {
	var msg string
	fields := gin.H{}

	if config.Config.App.Debug {
		msg = fmt.Sprintf("%v", err)
		fields["trace"] = string(debug.Stack())
	} else {
		msg = http.StatusText(http.StatusInternalServerError)
	}

	response(c, unknownErr, msg, nil, fields)
}

func ok(c *gin.Context, data interface{}, msg string) {
	response(c, statusOk, msg, data, nil)
}

func fail(c *gin.Context, code int, msg string) {
	response(c, code, msg, nil, nil)
}

func response(c *gin.Context, code int, msg string, data interface{}, fields gin.H) {
	resp := gin.H{
		"code":                 code,
		"msg":                  msg,
		"data":                 data,
		app_const.RequestIdKey: c.GetString(app_const.RequestIdKey),
	}

	for k, v := range fields {
		resp[k] = v
	}

	c.JSON(http.StatusOK, resp)
}
