package zaplog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// set the log format
func getEncoder() zapcore.Encoder {
	// return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// set the log file path
func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("/Users/yuan.ding/Desktop/code/entry_task/tcpserver/log/log.log")
	return zapcore.AddSync(file)
}

var Logger *zap.Logger
var Atom zap.AtomicLevel

func InitLogger() *zap.Logger {
	Atom = zap.NewAtomicLevel()
	encoder := getEncoder()
	writerSyncer := getLogWriter()
	consoleInfos := zapcore.Lock(os.Stdout)
	// 分别输出至 ./log 与 os.stdout
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writerSyncer, Atom),
		zapcore.NewCore(encoder, consoleInfos, Atom),
	)
	Logger = zap.New(core, zap.AddCaller())
	Atom.SetLevel(zap.ErrorLevel)
	return Logger
}

func init() {
	InitLogger()
}