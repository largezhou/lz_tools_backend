package api

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/largezhou/lz_tools_backend/app/config"
	"github.com/largezhou/lz_tools_backend/app/helper"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"github.com/largezhou/lz_tools_backend/app/model/user_model"
	"github.com/largezhou/lz_tools_backend/app/wechat"
	"go.uber.org/zap"
	"time"
)

var cfg = config.Config.App
var identityKey = user_model.GetIdentityKey()
var jwtMiddleware *jwt.GinJWTMiddleware

func getJwtMiddleware() *jwt.GinJWTMiddleware {
	if jwtMiddleware != nil {
		return jwtMiddleware
	}

	helper.CheckAppKey()

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "api",
		Key:         []byte(cfg.Key),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if user, ok := data.(*user_model.User); ok {
				return jwt.MapClaims{
					identityKey: user.Uuid,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			uuid := claims[identityKey].(string)
			return user_model.FindByUuid(uuid)
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var dto loginDto
			if err := c.ShouldBind(&dto); err != nil {
				return nil, err
			}
			res, err := wechat.OaOauth.GetUserAccessToken(dto.Code)
			if err != nil {
				logger.Error("获取 access_token 失败", zap.Error(err))
				return nil, err
			}

			userInfo, err := wechat.OaOauth.GetUserInfo(res.AccessToken, res.OpenID, "")
			if err != nil {
				logger.Error("获取微信用户信息失败", zap.Error(err))
				return nil, err
			}

			user := &user_model.User{
				OpenId:   userInfo.OpenID,
				UnionId:  userInfo.Unionid,
				Avatar:   userInfo.HeadImgURL,
				Nickname: userInfo.Nickname,
			}
			if err := user_model.UpdateOrCreateUserByUserInfo(user); err != nil {
				return nil, err
			}

			return user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if user, ok := data.(*user_model.User); ok && user != nil {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			logger.Info("登录失败", zap.Int("code", code), zap.String("msg", message))
			msg := ""
			if cfg.Debug {
				msg = message
			} else {
				msg = "登录失败"
			}
			fail(c, authFail, msg)
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			ok(c, gin.H{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			}, "")
		},
		TokenLookup:   "header: Authorization, query: _token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		panic(err)
	}

	return authMiddleware
}
