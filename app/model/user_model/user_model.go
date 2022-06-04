package user_model

import (
	"context"
	"github.com/google/uuid"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"github.com/largezhou/lz_tools_backend/app/model"
	"go.uber.org/zap"
)

type User struct {
	model.Model
	Uuid     string
	OpenId   string
	UnionId  string
	Nickname string
	Avatar   string
}

// GetIdentityKey 获取授权时的查询用的 key
func GetIdentityKey() string {
	return "uuid"
}

// Create 创建用户
func Create(ctx context.Context, user *User) error {
	user.Uuid = uuid.NewString()

	result := model.DB.WithContext(ctx).Create(user)
	return result.Error
}

// FindByUuid 通过 UUID 查找用户
func FindByUuid(ctx context.Context, uuid string) *User {
	var user *User
	if result := model.DB.WithContext(ctx).First(&user, "uuid = ?", uuid); result.Error != nil {
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

	return Create(ctx, userInfo)
}
