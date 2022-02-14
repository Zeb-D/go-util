package todo

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

//测试下zap框架，日志接入
func TestSimple(t *testing.T) {
	zap.NewProductionConfig()
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err)
	}
	//{"level":"info","ts":1639395809.250334,"caller":"todo/zap_test.go:18","msg":"log info"}
	logger.Info("log info")

	//获取sugarLogger，进行格式输出
	sugarLogger := logger.Sugar()                //Sugared 由Logger调用Sugar()方法获得
	sugarLogger.Infof("sugar loger %+v", logger) //支持printf风格输出日志
	//通用参数，可以理解是额外的非业务参数
	sugarLogger = sugarLogger.With("traceId", "123456")
	sugarLogger.Infof("2-sugar loger %+v", logger)
	//发现上面这些输出格式不太友好
	// 构造EncoderConfig
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity", //日志级别打印的key
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, //caller 打印的是全路径
	}
	// 构造 Config
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		InitialFields:    map[string]interface{}{"MyName": "kainhuck"},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	log, err := config.Build()
	log.Info("2-log info")
	log = log.With(zap.Any("Tid", 12))
	log.Info("2-log info")
}
