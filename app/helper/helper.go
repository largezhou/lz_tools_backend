package helper

import (
	"context"
	"errors"
	"github.com/largezhou/lz_tools_backend/app/app_error"
	"github.com/largezhou/lz_tools_backend/app/config"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"gorm.io/gorm"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type ShutdownFunc func(ctx context.Context)

var shutdownFuncList []ShutdownFunc

// RegisterShutdownFunc 注册一个服务关闭时的回调函数
func RegisterShutdownFunc(f ShutdownFunc) {
	shutdownFuncList = append(shutdownFuncList, f)
}

// CallShutdownFunc 服务关闭时，执行所有回调函数
func CallShutdownFunc(ctx context.Context) {
	for _, f := range shutdownFuncList {
		f(ctx)
	}
}

// CheckAppKey 检查 app key
func CheckAppKey() {
	if len(config.Config.App.Key) < 32 {
		panic("APP_KEY 长度至少为 32 位")
	}
}

// ModelNotFound 处理模型未找到
func ModelNotFound(err error, msg string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return app_error.New(msg).SetCode(app_error.ResourceNotFound)
	} else {
		return err
	}
}

// GetQrcodeFromFile 从上传文件中识别二维码
func GetQrcodeFromFile(ctx context.Context, file io.Reader) (string, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	qrReader := qrcode.NewQRCodeReader()
	qrRes, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}
	if qrRes == nil {
		return "", app_error.New("无法识别有效的二维码")
	}

	return qrRes.String(), nil
}
