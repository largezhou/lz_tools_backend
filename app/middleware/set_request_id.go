package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/largezhou/lz_tools_backend/app/app_const"
)

func SetRequestId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestId := ctx.GetHeader(app_const.RequestIdHeaderKey)
		if requestId == "" {
			requestId = uuid.NewString()
		}
		ctx.Set(app_const.RequestIdKey, requestId)
		ctx.Header(app_const.RequestIdHeaderKey, requestId)

		ctx.Next()
	}
}
