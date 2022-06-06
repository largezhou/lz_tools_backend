package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/largezhou/lz_tools_backend/app/app_error"
	"github.com/largezhou/lz_tools_backend/app/dto/code_dto"
	"github.com/largezhou/lz_tools_backend/app/helper"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"github.com/largezhou/lz_tools_backend/app/model"
	"github.com/largezhou/lz_tools_backend/app/model/code_model"
	tRedis "github.com/largezhou/lz_tools_backend/app/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"strings"
)

type CodeService struct {
}

var redisService = tRedis.Client

func NewCodeService() *CodeService {
	return &CodeService{}
}

func (cs *CodeService) GetCodeList(
	ctx context.Context,
	userId uint,
	dto code_dto.GetCodeListDto,
) ([]*code_model.Code, error) {
	codeList, _ := code_model.GetCodeByUserId(ctx, userId)
	return codeList, nil
}

func (cs CodeService) SaveCode(ctx context.Context, userId uint, dto code_dto.SaveCodeDto) error {
	if dto.Id == 0 {
		return cs.CreateCode(ctx, userId, dto)
	} else {
		return cs.UpdateCode(ctx, userId, dto)
	}
}

func (cs *CodeService) CreateCode(ctx context.Context, userId uint, dto code_dto.SaveCodeDto) error {
	db := model.DB.WithContext(ctx)
	var code *code_model.Code
	if dto.CopyFromId > 0 {
		var copyFrom *code_model.Code
		if result := db.Where("id = ?", dto.CopyFromId).
			Where("share = ?", true).
			First(&copyFrom); result.Error != nil {
			return helper.ModelNotFound(result.Error, "无法复制场所码")
		}

		code = &code_model.Code{
			CopyFromId: dto.CopyFromId,
			UserId:     userId,
			Name:       copyFrom.Name,
			Lng:        copyFrom.Lng,
			Lat:        copyFrom.Lat,
			Link:       copyFrom.Link,
			Share:      false,
		}
	} else {
		dto.Name = strings.TrimSpace(dto.Name)
		if dto.Name == "" || dto.Lng <= 0 || dto.Lat <= 0 {
			return app_error.New("场所名或经纬度不能为空")
		}
		if dto.File == nil {
			return app_error.New("需要上传场所码文件")
		}

		link, err := cs.getQrcodeFromUploadedFile(ctx, dto.File)
		if err != nil {
			return err
		}

		code = &code_model.Code{
			CopyFromId: dto.CopyFromId,
			UserId:     userId,
			Name:       dto.Name,
			Lng:        dto.Lng,
			Lat:        dto.Lat,
			Link:       link,
			Share:      dto.Share,
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&code); result.Error != nil {
			return result.Error
		}

		return cs.updateRedisGeo(ctx, code)
	})

	return err
}

func (cs *CodeService) updateRedisGeo(ctx context.Context, code *code_model.Code) error {
	_, err := redisService.GeoAdd(ctx, "code_geo", &redis.GeoLocation{
		Name:      fmt.Sprintf("%d", code.Id),
		Longitude: code.Lng,
		Latitude:  code.Lat,
	}).Result()
	if err != nil {
		logger.Error(ctx, "redis 坐标添加失败", zap.Error(err))
	}
	return err
}

func (cs *CodeService) UpdateCode(ctx context.Context, userId uint, dto code_dto.SaveCodeDto) error {
	db := model.DB.WithContext(ctx)
	var code *code_model.Code
	if result := db.Where("id = ?", dto.Id).
		Where("user_id = ?", userId).
		First(&code); result.Error != nil {
		return helper.ModelNotFound(result.Error, "场所码不存在")
	}

	if dto.Name == "" || dto.Lng <= 0 || dto.Lat <= 0 {
		return app_error.New("场所名或经纬度不能为空")
	}

	code.Name = dto.Name
	code.Share = dto.Share
	code.Often = dto.Often
	code.Lng = dto.Lng
	code.Lat = dto.Lat

	if dto.File != nil {
		link, err := cs.getQrcodeFromUploadedFile(ctx, dto.File)
		if err != nil {
			return err
		}
		code.Link = link
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if result := tx.Updates(&code); result.Error != nil {
			return result.Error
		}

		return cs.updateRedisGeo(ctx, code)
	})

	return err
}

func (cs *CodeService) getQrcodeFromUploadedFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	link, err := helper.GetQrcodeFromFile(ctx, file)
	if err != nil {
		logger.Error(ctx, "识别二维码失败", zap.Error(err))
		return "", app_error.New("识别二维码失败")
	}

	return link, nil
}

func (cs CodeService) DeleteCode(ctx context.Context, userId uint, codeId uint) error {
	code, err := code_model.GetCodeByIdAndUserId(ctx, codeId, userId)
	if err != nil {
		return helper.ModelNotFound(err, "场所码不存在")
	}

	err = model.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if res := tx.Delete(&code); res.Error != nil {
			return res.Error
		}
		if _, err := redisService.ZRem(ctx, "code_geo", code.Id).Result(); err != nil {
			return err
		}
		return nil
	})

	return err
}