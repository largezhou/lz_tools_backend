package api

import (
	"github.com/gin-gonic/gin"
	"github.com/largezhou/lz_tools_backend/app/app_error"
	"github.com/largezhou/lz_tools_backend/app/dto"
	"github.com/largezhou/lz_tools_backend/app/dto/code_dto"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"github.com/largezhou/lz_tools_backend/app/service"
	"go.uber.org/zap"
	"path/filepath"
)

type CodeController struct {
	codeService *service.CodeService
}

func NewCodeController() *CodeController {
	return &CodeController{
		codeService: service.NewCodeService(),
	}
}

func (cc *CodeController) GetCodeList(ctx *gin.Context) {
	var req code_dto.GetCodeListDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		failWithError(ctx, err)
		return
	}

	user, _ := getAuthUser(ctx)
	codeList, _ := cc.codeService.GetCodeList(ctx, user.Id, req)

	ok(ctx, codeList, "")
}

func (cc CodeController) SaveCode(ctx *gin.Context) {
	var req code_dto.SaveCodeDto
	if err := ctx.ShouldBind(&req); err != nil {
		failWithError(ctx, err)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Info(ctx, "无上传文件", zap.Error(err))
	}
	allowExt := map[string]bool{
		".png":  true,
		".jpeg": true,
		".jpg":  true,
	}
	if file != nil {
		if _, ok := allowExt[filepath.Ext(file.Filename)]; !ok {
			fail(ctx, app_error.InvalidParameter, "图片必须是 png，jpeg 或者 jpg 格式")
			return
		}
		req.File = file
	}

	user, _ := getAuthUser(ctx)
	if err := cc.codeService.SaveCode(ctx, user.Id, req); err != nil {
		failWithError(ctx, err)
		return
	}

	ok(ctx, nil, "")
}

func (cc CodeController) DeleteCode(ctx *gin.Context) {
	var req dto.IdDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		failWithError(ctx, err)
		return
	}

	user, _ := getAuthUser(ctx)
	if err := cc.codeService.DeleteCode(ctx, user.Id, req.Id); err != nil {
		failWithError(ctx, err)
		return
	}

	ok(ctx, nil, "")
}
