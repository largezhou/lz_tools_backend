package model

import (
	"context"
	"github.com/largezhou/lz_tools_backend/app/logger"
	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

type SqlRecorderLogger struct {
	gormLogger.Interface
}

func (l SqlRecorderLogger) Trace(
	ctx context.Context,
	_ time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	sql, rows := fc()
	logger.Info(ctx,"sql", zap.String("sql", sql), zap.Int64("rows", rows), zap.Error(err))
}
