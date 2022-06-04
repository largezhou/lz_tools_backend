package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/largezhou/lz_tools_backend/app/model"
	"github.com/largezhou/lz_tools_backend/app/model/user_model"
	"github.com/largezhou/lz_tools_backend/app/wechat"
)

func getCode(ctx *gin.Context) {
	model.DB.WithContext(ctx).First(&user_model.User{})

	ok(ctx, "https://www.baidu.com", "")
}

func getCodeList(ctx *gin.Context) {
	ok(ctx, nil, "")
}

func getWechatAuthUrl(ctx *gin.Context) {
	var req getWechatAuthUrlDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if errors.As(err, &validator.ValidationErrors{}) {
			fail(ctx, invalidParameter, err.Error())
		} else {
			fail(ctx, badRequest, "")
		}

		return
	}

	url, _ := wechat.OfficialAccount.GetOauth().GetRedirectURL(
		req.Redirect,
		"snsapi_userinfo",
		"",
	)

	ok(ctx, gin.H{"url": url}, "")
}
