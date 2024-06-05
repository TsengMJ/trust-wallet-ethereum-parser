package logger

import (
	"errors"
	"ethereum-parser/config"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLog() error {
	logConfig := config.Config.Log

	if logConfig.Path == "" {
		return errors.New("Log path is empty")
	}

	if logConfig.ErrorLog == "" {
		return errors.New("Error log is empty")
	}

	if logConfig.WarnLog == "" {
		return errors.New("Warn log is empty")
	}

	if logConfig.InfoLog == "" {
		return errors.New("Info log is empty")
	}

	if logConfig.DebugLog == "" {
		return errors.New("Debug log is empty")
	}

	err := os.MkdirAll(logConfig.Path, os.ModePerm)
	if err != nil {
		return errors.New("Error creating log directory, " + err.Error())
	}

	errorLogPath := filepath.Join(logConfig.Path, logConfig.ErrorLog)
	errorCore, err := getZapCore(errorLogPath, zapcore.ErrorLevel)
	if err != nil {
		return errors.New("Error creating error core, " + err.Error())
	}

	warnLogPath := filepath.Join(logConfig.Path, logConfig.WarnLog)
	warnCore, err := getZapCore(warnLogPath, zapcore.WarnLevel)
	if err != nil {
		return errors.New("Error creating warn core, " + err.Error())
	}

	infoLogPath := filepath.Join(logConfig.Path, logConfig.InfoLog)

	infoCore, err := getZapCore(infoLogPath, zapcore.InfoLevel)
	if err != nil {
		return errors.New("Error creating info core, " + err.Error())
	}

	debugLogPath := filepath.Join(logConfig.Path, logConfig.DebugLog)
	debugCore, err := getZapCore(debugLogPath, zapcore.DebugLevel)
	if err != nil {
		return errors.New("Error creating debug core, " + err.Error())
	}

	core := zapcore.NewTee([]zapcore.Core{
		errorCore, warnCore, infoCore, debugCore,
	}...)

	Logger = zap.New(core, zap.AddCaller())

	return nil
}

func getZapCore(logFile string, logLevel zapcore.Level) (zapcore.Core, error) {
	logConfig := config.Config.Log

	logWriter, err := rotatelogs.New(
		logFile+".%Y-%m-%d",
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(time.Duration(logConfig.MaxAge)*time.Hour),
	)
	if err != nil {
		return nil, errors.New("Error creating log writer, " + err.Error())
	}

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	errorPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == logLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), errorPriority),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logWriter), errorPriority),
	)

	return core, nil
}
