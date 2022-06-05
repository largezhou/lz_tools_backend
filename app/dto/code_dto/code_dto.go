package code_dto

import (
	"mime/multipart"
)

type GetWechatAuthUrlDto struct {
	Redirect string `binding:"required"`
}

type LoginDto struct {
	Code string `binding:"required"`
}

type GetCodeListDto struct {
	Lng float64
	Lat float64
}

type SaveCodeDto struct {
	Id         uint    `form:"id"`
	Name       string  `form:"name" binding:"omitempty,gt=0"`
	Lng        float64 `form:"lng" binding:"omitempty,gt=0"`
	Lat        float64 `form:"lat" binding:"omitempty,gt=0"`
	CopyFromId uint    `form:"copyFromId"`
	Often      bool    `form:"often"`
	Share      bool    `form:"share"`
	File       *multipart.FileHeader
}
