package user_model

import (
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

var db = model.DB

// GetIdentityKey 获取授权时的查询用的 key
func GetIdentityKey() string {
	return "uuid"
}

// Create 创建用户
func Create(user *User) error {
	user.Uuid = uuid.NewString()

	result := db.Create(user)
	return result.Error
}

// FindByUuid 通过 UUID 查找用户
func FindByUuid(uuid string) *User {
	var user *User
	if result := db.First(&user, "uuid = ?", uuid); result.Error != nil {
		return nil
	}
	return user
}

// UpdateOrCreateUserByUserInfo 通过用户信息中的 union_id 查找，或者创建用户
func UpdateOrCreateUserByUserInfo(userInfo *User) error {
	var user *User
	// 查找到用户，则更新信息，否则新建用户
	if result := db.First(&user, "union_id = ?", userInfo.UnionId); result.Error == nil {
		if updateResult := db.Where("union_id = ?", userInfo.UnionId).Updates(userInfo); updateResult.Error != nil {
			logger.Warn("用户更新失败", zap.Error(updateResult.Error))
		}
		*userInfo = *user
		return nil
	}

	return Create(userInfo)
}
