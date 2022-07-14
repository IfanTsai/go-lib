package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/IfanTsai/go-lib/logger"
	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	infoStr = "%s\n[info] "
	warnStr = "%s\n[warn] "
	errStr  = "%s\n[error] "
)

type Logger struct {
	gormLogger.Config
	zapLogger *zap.Logger
}

func newLogger(filename string) *Logger {
	return &Logger{
		zapLogger: logger.NewJSONLogger(
			logger.WithDisableConsole(),
			logger.WithFileRotationP(filename),
		),
	}
}

func (l *Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.LogLevel = level

	return l
}

func (l *Logger) Info(ctx context.Context, msg string, values ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		l.zapLogger.Info(fmt.Sprintf(infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, values...)...))
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, values ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		l.zapLogger.Warn(fmt.Sprintf(warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, values...)...))
	}
}

func (l *Logger) Error(ctx context.Context, msg string, values ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		l.zapLogger.Error(fmt.Sprintf(errStr+msg, append([]interface{}{utils.FileWithLineNum()}, values...)...))
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	sql, rowsAffected := fc()
	elapsed := time.Since(begin)
	duration := float64(elapsed.Nanoseconds()) / float64(time.Millisecond)
	if err != nil {
		l.zapLogger.Error(
			"execute_sql",
			zap.String("source", utils.FileWithLineNum()),
			zap.Float64("duration", duration),
			zap.Int64("rows_affected", rowsAffected),
			zap.String("sql", sql),
			zap.Error(err),
		)
	} else {
		l.zapLogger.Info(
			"execute_sql",
			zap.String("source", utils.FileWithLineNum()),
			zap.Float64("duration", duration),
			zap.Int64("rows_affected", rowsAffected),
			zap.String("sql", sql),
		)
	}
}
