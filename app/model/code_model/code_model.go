package code_model

import (
	"context"
	"github.com/largezhou/lz_tools_backend/app/model"
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

func GetCodeByUserId(ctx context.Context, userId uint) ([]*Code, error) {
	var codeList []*Code
	result := model.DB.WithContext(ctx).
		Where("user_id = ?", userId).
		Order("often desc, times desc").
		Find(&codeList)
	return codeList, result.Error
}

func GetCodeByIdAndUserId(ctx context.Context, id uint, userId uint) (*Code, error) {
	var code *Code
	res := model.DB.WithContext(ctx).Where("id = ?", id).Where("user_id = ?", userId).First(&code)
	return code, res.Error
}
