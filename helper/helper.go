package helper

import (
	"os"
	"reflect"
	"regexp"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Var
var (
	pattern *regexp.Regexp = regexp.MustCompile(`(?im)(\s*)([^\s]+)`)
	logger  *zap.SugaredLogger
)

// InitLogger : Init logger with debug option
func InitLogger(isDebug bool) *zap.SugaredLogger {
	loggerMgr := initZapLog(isDebug)
	zap.ReplaceGlobals(loggerMgr)
	defer func() {
		// flushes buffer, if any
		_ = zap.S().Sync()
	}()
	return zap.S()
}

// GetLogger - call InitLogger first, otherwise, return a default logger
func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		logger = zap.S()
	}
	return logger
}

// Init Zap Logger
func initZapLog(isDebug bool) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	if isDebug {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		config.EncoderConfig.LevelKey = ""
		config.EncoderConfig.TimeKey = ""
		config.EncoderConfig.CallerKey = ""
	}
	logger, _ := config.Build()
	return logger
}

// HandleError : Handle an error without exiting (just logging)
func HandleError(e error) {
	if e != nil {
		GetLogger().Error(e)
	}
}

// HandleErrorExit : Handle an error and exit
func HandleErrorExit(e error) {
	if e != nil {
		GetLogger().Fatal(e)
		os.Exit(1)
	}
}

// Ternary : Execute ternery condition with condition, return resultTrue if condition is true, restultFalse otherwise
func Ternary(condition bool, resultTrue, restultFalse interface{}) interface{} {
	if condition {
		if reflect.TypeOf(resultTrue).Kind() == reflect.Func {
			return resultTrue.(func() interface{})()
		}
		return resultTrue
	}
	if reflect.TypeOf(restultFalse).Kind() == reflect.Func {
		return restultFalse.(func() interface{})()
	}
	return restultFalse
}

// MaskPassword : will mask any password within the given string and return it.
// Passwords are identified as following the phrase "password" or
// "pass" and an equal(=) or colon(:) sign, and are replaced with "****"
func MaskPassword(toMask string) string {
	length := len([]rune(toMask))
	return pattern.ReplaceAllString(toMask, strings.Repeat("*", length))
}
