package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
}

func InitDefault(logLevel zapcore.Level) {
	stdout := zapcore.AddSync(os.Stdout)
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(productionCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, logLevel),
		zapcore.NewCore(fileEncoder, file, logLevel),
	)

	Init(core)
}
