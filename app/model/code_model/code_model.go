package code_model

import (
	"context"
	"github.com/largezhou/lz_tools_backend/app/dto/code_dto"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"github.com/largezhou/lz_tools_backend/app/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Code struct {
	model.Model
	CopyFromId uint    `json:"copyFromId"`
	UserId     uint    `json:"userId"`
	Name       string  `json:"name"`
	Lng        float64 `json:"lng"`
	Lat        float64 `json:"lat"`
	Link       string  `json:"link"`
	Times      uint    `json:"times"`
	Often      bool    `json:"often"`
	Share      bool    `json:"share"`
}

func GetCodeByUserId(ctx context.Context, userId uint) ([]*code_dto.CodeListDto, error) {
	var codeList []*code_dto.CodeListDto
	result := model.DB.WithContext(ctx).
		Model(&Code{}).
		Where("user_id = ?", userId).
		Order("often desc, times desc, id desc").
		Find(&codeList)
	return codeList, result.Error
}

func GetCodeByIdAndUserId(ctx context.Context, id uint, userId uint) (*Code, error) {
	var code *Code
	res := model.DB.WithContext(ctx).Where("id = ?", id).Where("user_id = ?", userId).First(&code)
	return code, res.Error
}

func UpdateTimes(ctx context.Context, id uint) {
	if res := model.DB.WithContext(ctx).
		Model(&Code{}).
		Where("id = ?", id).
		Update("times", gorm.Expr("times + 1")); res.Error != nil {
		logger.Info(ctx, "更新使用次数失败", zap.Error(res.Error))
	}
}
