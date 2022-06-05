package wechat

import (
	"context"
	"github.com/largezhou/lz_tools_backend/app/config"
	wechatPkg "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	oaPkg "github.com/silenceper/wechat/v2/officialaccount"
	oaConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
)

var Service *wechatPkg.Wechat

// OfficialAccount 微信公众号
var OfficialAccount *oaPkg.OfficialAccount

// OaOauth 公众号 oauth
var OaOauth *oauth.Oauth

var cfg = config.Config

func init() {
	Service = wechatPkg.NewWechat()
	Service.SetCache(cache.NewRedis(context.Background(), &cache.RedisOpts{
		Host:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		Database: cfg.Redis.Db,
	}))

	OfficialAccount = Service.GetOfficialAccount(&oaConfig.Config{
		AppID:     cfg.Wechat.OfficialAccount.AppId,
		AppSecret: cfg.Wechat.OfficialAccount.AppSecret,
		Token:     cfg.Wechat.OfficialAccount.Token,
	})

	OaOauth = OfficialAccount.GetOauth()
}
