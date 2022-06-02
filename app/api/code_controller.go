package api

import (
	"github.com/gin-gonic/gin"
)

func GetCode(c *gin.Context) {
	ok(c, "https://www.baidu.com", "")
}

func GetCodeList(c *gin.Context) {
	fail(c, authFail, "")
}
