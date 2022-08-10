package main

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger = initLog(filepath.Join(getLogDir(), "log.json"))
var defaultLogLevel = zapcore.DebugLevel

var timeLayout = "2006-01-02T15:04:05.000Z"

func getLogDir() string {
	return os.Getenv("LOG_DIR")
}

func initLog(filePath string) *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	//config.EncodeTime = zapcore.ISO8601TimeEncoder
	//config.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncodeTime = zapcore.TimeEncoderOfLayout(timeLayout)
	zCores := []zapcore.Core{
		getFileCore(filePath, config),
		getStdoutCore(config),
	}
	core := zapcore.NewTee(zCores...)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func getFileCore(filePath string, config zapcore.EncoderConfig) zapcore.Core {
	fileEncoder := zapcore.NewJSONEncoder(config)
	logFile, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // not handling error as of now
	defer logFile.Close()
	lumlog := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     2,
	}
	return zapcore.NewCore(fileEncoder, zapcore.AddSync(lumlog), defaultLogLevel) // writes to file
	// return zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), defaultLogLevel)
}

func getStdoutCore(config zapcore.EncoderConfig) zapcore.Core {
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	return zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel) // writes to stdout
}

func main() {
	uid := "j1ae77e5-915l-13e9-bc42-556nf7864d64"
	uidZapfield := func(u string) zapcore.Field {
		return zap.String("messageid", u)
	}
	statusZapField := func(s int) zapcore.Field {
		return zap.Int("statuscode", s)
	}
	methodZapField := func(m string) zapcore.Field {
		return zap.String("method", m)
	}
	for i := 0; i < 100; i++ {
		logger.Info("Hello World", zap.String("reason", "Greeting from unidentifed species"), statusZapField(200), methodZapField("PUT"), uidZapfield(uid))
		time.Sleep(250 * time.Millisecond)
		logger.Error("End of World", zap.String("reason", "Invasion of unidentifed species"), statusZapField(500), methodZapField("PUT"), uidZapfield(uid))
		time.Sleep(250 * time.Millisecond)
	}
}
