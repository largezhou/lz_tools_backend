package model

import (
	"github.com/largezhou/gin_starter/app/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
	"strings"
	"time"
)

var c = config.Config.Mysql
var appConfig = config.Config.App
var DB *gorm.DB

type Model struct {
	Id         uint      `gorm:"primaryKey"`
	CreateTime time.Time `gorm:"type:datetime;autoCreateTime"`
	UpdateTime time.Time `gorm:"type:datetime;autoUpdateTime "`
}

func init() {
	dsn := c.Dsn
	if !strings.Contains(dsn, "loc=") {
		dsn += "&loc=" + url.QueryEscape(appConfig.Timezone)
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: &SqlRecorderLogger{},
	})
	if err != nil {
		panic(err)
	}
}
