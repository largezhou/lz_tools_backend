package model

import (
	"github.com/largezhou/lz_tools_backend/app/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/url"
	"strings"
	"time"
)

var c = config.Config.Mysql
var appConfig = config.Config.App
var DB *gorm.DB

type Model struct {
	Id         uint      `gorm:"primaryKey" json:"id"`
	CreateTime time.Time `gorm:"type:datetime;autoCreateTime;not null" json:"createTime"`
	UpdateTime time.Time `gorm:"type:datetime;autoUpdateTime;not null" json:"updateTime"`
}

func init() {
	dsn := c.Dsn
	if !strings.Contains(dsn, "loc=") {
		dsn += "&loc=" + url.QueryEscape(appConfig.Timezone)
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: &SqlRecorderLogger{},
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
}
