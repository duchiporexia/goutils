package xlog

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/rs/zerolog"

	gl "gorm.io/gorm/logger"
)

type Logger struct {
}

func NewGormLogger() *Logger {
	return &Logger{}
}

func (l Logger) LogMode(gl.LogLevel) gl.Interface {
	return l
}

func (l Logger) Error(ctx context.Context, msg string, opts ...interface{}) {
	if needGormLog {
		logger.Error().Msg(fmt.Sprintf(msg, opts...))
	}
}

func (l Logger) Warn(ctx context.Context, msg string, opts ...interface{}) {
	if needGormLog {
		logger.Warn().Msg(fmt.Sprintf(msg, opts...))
	}
}

func (l Logger) Info(ctx context.Context, msg string, opts ...interface{}) {
	if needGormLog {
		logger.Info().Msg(fmt.Sprintf(msg, opts...))
	}
}

func (l Logger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	if !needGormLog {
		return
	}

	zl := logger
	var event *zerolog.Event

	if err != nil && err != gorm.ErrRecordNotFound {
		event = zl.Error().Err(err)
	} else {
		event = zl.Info()
	}

	event = event.Caller(3)

	sql, rows := f()
	if rows > -1 {
		event.Msg(fmt.Sprintf("[rows:%d] %s", rows, sql))
	} else {
		event.Msg(sql)
	}

	return
}
