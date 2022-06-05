package user_model

import (
	"context"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"github.com/largezhou/lz_tools_backend/app/model"
	"go.uber.org/zap"
)

type User struct {
	model.Model
	OpenId   string `json:"openId"`
	UnionId  string `json:"unionId"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
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
func UpdateOrCreateUserByUserInfo(ctx context.Context, userInfo *User) error {
	var user *User
	db := model.DB.WithContext(ctx)
	// 查找到用户，则更新信息，否则新建用户
	if result := db.First(&user, "union_id = ?", userInfo.UnionId); result.Error == nil {
		if updateResult := db.Where("union_id = ?", userInfo.UnionId).Updates(userInfo); updateResult.Error != nil {
			logger.Warn(ctx, "用户更新失败", zap.Error(updateResult.Error))
		}
		*userInfo = *user
		return nil
	}

	return model.DB.WithContext(ctx).Create(user).Error
}
