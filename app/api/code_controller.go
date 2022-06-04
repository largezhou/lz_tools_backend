package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/largezhou/lz_tools_backend/app/wechat"
)

func getCode(c *gin.Context) {
	ok(c, "https://www.baidu.com", "")
}

func getCodeList(c *gin.Context) {
	ok(c, nil, "")
}

func getWechatAuthUrl(c *gin.Context) {
	var req getWechatAuthUrlDto
	if err := c.ShouldBindJSON(&req); err != nil {
		if errors.As(err, &validator.ValidationErrors{}) {
			fail(c, invalidParameter, err.Error())
		} else {
			fail(c, badRequest, "")
		}

		return
	}

	url, _ := wechat.OfficialAccount.GetOauth().GetRedirectURL(
		req.Redirect,
		"snsapi_userinfo",
		"",
	)

	ok(c, gin.H{"url": url}, "")
}
