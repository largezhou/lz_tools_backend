package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 简单跨域中间件
func Cors(r gin.IRoutes) gin.HandlerFunc {
	corsFunc := func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
	r.Use(corsFunc).OPTIONS("/*ignore", func(c *gin.Context) {
		return
	})

	return corsFunc
}
