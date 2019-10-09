package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	FormatEnumJSON    FormatEnum = "json"
	FormatEnumConsole FormatEnum = "console"

	DefaultRotatePeriodSecond int64 = 60 * 60 * 24 // 一天
)

var (
	LevelEnumDebug = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	LevelEnumInfo  = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	LevelEnumWarn  = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	LevelEnumError = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
)

//默认开发环境
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		localTime:          true,
		multiWrite:         true,
		maxSize:            10 * 1024 * 1024, //
		rotatePeriodSecond: DefaultRotatePeriodSecond,
		level:              LevelEnumDebug,
		appName:            "tmp",
		dir:                "logs",
		format:             FormatEnumConsole,
	}
}

//线上环境
func ProdLogConfig(appName string) *LogConfig {
	c := DefaultLogConfig()
	c.multiWrite = false
	c.level = LevelEnumInfo
	c.appName = appName
	c.format = FormatEnumJSON
	return c
}

type LogConfig struct {
	localTime bool
	// MultiWrite 如果为true，日志会同时写入文件和终端，默认为true
	multiWrite bool
	// 日志级别，默认为debug
	level LevelEnum
	// AppName 默认为tmp
	appName string
	// Dir 日志文件存放目录，默认为logs
	dir string
	// Format默认为console，方便调试
	format FormatEnum
	// hook
	Hooks []LogHook
	// 单个log文件最大size，默认为100M
	maxSize int
	// 日志按时间分割，默认为一天
	rotatePeriodSecond int64
}

type LevelEnum = zap.AtomicLevel

type FormatEnum string

type FormatTime = zapcore.TimeEncoder
