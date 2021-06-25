package xlog

import (
	"fmt"
	utilscommon "github.com/duchiporexia/goutils/internal"
	"github.com/duchiporexia/goutils/xutils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"strconv"
	"strings"
	"time"
)

var logger zerolog.Logger
var Ctxt zerolog.Context

var needGormLog = true

func Init(env string) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		return xutils.RelativePath(file) + ":" + strconv.Itoa(line)
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	hasInit := false
	for _, v := range []string{"dev", "stg", "uat", "test", "local"} {
		if strings.HasPrefix(env, v) {
			logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
			hasInit = true
		}
	}
	if !hasInit {
		needGormLog = false
		logger = zerolog.New(os.Stdout)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	Ctxt = logger.With().Timestamp()

	logger = Ctxt.Logger()
}

func DebugE() *zerolog.Event {
	return logger.Debug().Caller(1)
}

func InfoE() *zerolog.Event {
	return logger.Info().Caller(1)
}
func Info(msg string) {
	logger.Info().Caller(1).Msg(msg)
}
func InfoF(format string, a ...interface{}) {
	logger.Info().Caller(1).Msg(fmt.Sprintf(format, a...))
}

func WarnE() *zerolog.Event {
	return logger.Warn().Caller(1)
}
func Warn(msg string) {
	logger.Warn().Caller(1).Msg(msg)
}
func WarnF(format string, a ...interface{}) {
	logger.Warn().Caller(1).Msg(fmt.Sprintf(format, a...))
}

func ErrorE(skip int) *zerolog.Event {
	return logger.Error().Caller(skip)
}
func Error(err error) {
	msg := fmt.Sprintf("%+v", utilscommon.WithStack(err))
	logger.Error().Caller(1).Msg(msg)
}
func ErrorF(format string, a ...interface{}) {
	logger.Error().Caller(1).Msg(fmt.Sprintf(format, a...))
}

func FatalE() *zerolog.Event {
	return logger.Fatal().Caller(1)
}
func Fatal(msg string) {
	logger.Fatal().Caller(1).Msg(msg)
}
