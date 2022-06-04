package wechat

import (
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

func init() {
	c := config.Config.Wechat.OfficialAccount
	Service = wechatPkg.NewWechat()
	OfficialAccount = Service.GetOfficialAccount(&oaConfig.Config{
		AppID:     c.AppId,
		AppSecret: c.AppSecret,
		Token:     c.Token,
		Cache:     cache.NewMemory(),
	})

	OaOauth = OfficialAccount.GetOauth()
}
