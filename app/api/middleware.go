package api

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/largezhou/lz_tools_backend/app/app_const"
	"github.com/largezhou/lz_tools_backend/app/app_error"
	"github.com/largezhou/lz_tools_backend/app/config"
	"github.com/largezhou/lz_tools_backend/app/dto/code_dto"
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
		Timeout:     24 * 30 * time.Hour,
		MaxRefresh:  24 * 60 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data any) jwt.MapClaims {
			if user, ok := data.(*user_model.User); ok {
				return jwt.MapClaims{
					identityKey: user.Id,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx *gin.Context) any {
			claims := jwt.ExtractClaims(ctx)
			id := uint(claims[identityKey].(float64))
			return user_model.FindById(ctx, id)
		},
		Authenticator: func(ctx *gin.Context) (any, error) {
			var dto code_dto.LoginDto
			if err := ctx.ShouldBind(&dto); err != nil {
				return nil, err
			}
			res, err := wechat.OaOauth.GetUserAccessToken(dto.Code)
			if err != nil {
				logger.Error(ctx, "获取 access_token 失败", zap.Error(err))
				return nil, err
			}

			userInfo, err := wechat.OaOauth.GetUserInfo(res.AccessToken, res.OpenID, "")
			if err != nil {
				logger.Error(ctx, "获取微信用户信息失败", zap.Error(err))
				return nil, err
			}

			user := &user_model.User{
				OpenId:   userInfo.OpenID,
				UnionId:  userInfo.Unionid,
				Avatar:   userInfo.HeadImgURL,
				Nickname: userInfo.Nickname,
			}
			return user_model.UpdateOrCreateUserByUserInfo(ctx, user)
		},
		Authorizator: func(data any, ctx *gin.Context) bool {
			if user, ok := data.(*user_model.User); ok && user != nil {
				ctx.Set(app_const.AuthUserKey, user)
				return true
			}

			return false
		},
		Unauthorized: func(ctx *gin.Context, code int, message string) {
			logger.Info(ctx, "登录失败", zap.Int("code", code), zap.String("msg", message))
			msg := ""
			if cfg.Debug {
				msg = message
			} else {
				msg = "登录失败"
			}
			fail(ctx, app_error.AuthFail, msg)
		},
		LoginResponse: func(ctx *gin.Context, code int, token string, expire time.Time) {
			ok(ctx, gin.H{
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
