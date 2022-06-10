package code_dto

import (
	"mime/multipart"
	"time"
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
	Id    uint    `form:"id"`
	Name  string  `form:"name" binding:"required,gt=0"`
	Lng   float64 `form:"lng" binding:"required,gt=0"`
	Lat   float64 `form:"lat" binding:"required,gt=0"`
	Often bool    `form:"often"`
	File  *multipart.FileHeader
}

type CodeListDto struct {
	Id         uint      `json:"id"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	UserId     uint      `json:"userId"`
	Name       string    `json:"name"`
	Lng        float64   `json:"lng"`
	Lat        float64   `json:"lat"`
	Link       string    `json:"link"`
	Times      uint      `json:"times"`
	Often      bool      `json:"often"`
	Dist       float64   `json:"dist"`
}
