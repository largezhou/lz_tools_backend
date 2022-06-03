package user_model

import (
	"github.com/largezhou/gin_starter/app/model"
)

type User struct {
	model.Model
	UUID     string `gorm:"uniqueIndex;not null"`
	UnionId  string `gorm:"uniqueIndex;not null"`
	Nickname string
	Avatar   string
}

func (User) TableName() string {
	return "user"
}
