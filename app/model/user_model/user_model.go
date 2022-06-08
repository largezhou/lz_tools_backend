package user_model

import (
	"context"
	"github.com/largezhou/lz_tools_backend/app/model"
	"gorm.io/gorm"
)

type User struct {
	model.Model
	Username string `json:"nickname"`
}

// GetIdentityKey 获取授权时的查询用的 key
func GetIdentityKey() string {
	return "id"
}

// FindById 通过 UUID 查找用户
func FindById(ctx context.Context, id uint) *User {
	var user *User
	if result := model.DB.WithContext(ctx).First(&user, "id = ?", id); result.Error != nil {
		return nil
	}
	return user
}

// UpdateOrCreateUserByUserInfo 通过用户信息中的 union_id 查找，或者创建用户
func UpdateOrCreateUserByUserInfo(ctx context.Context, userInfo *User) (*User, error) {
	db := model.DB.WithContext(ctx)

	var result *gorm.DB
	result = db.Where("username = ?", userInfo.Username).Updates(userInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		db.First(&userInfo, "username = ?", userInfo.Username)
		return userInfo, nil
	}

	if result = db.Create(&userInfo); result.Error != nil {
		return nil, result.Error
	}

	return userInfo, result.Error
}
