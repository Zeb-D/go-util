package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var log *zap.Logger
var sugarLog *zap.SugaredLogger

//simple test println
func sayHello(appName string, name string) {
	println("%+v", appName, name)
}

//封装个zap日志打印
func zapInfoLog(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func SetGlobalLog(appName string, prod bool, opts ...OptionFunc) {
	if appName == "" {
		panic("appName must not empty")
	}
	var config *LogConfig
	if prod {
		config = ProdLogConfig(appName)
	} else {
		config = DefaultLogConfig()
		config.appName = appName
	}
	for _, v := range opts {
		v(config)
	}

	log = NewLog(config)
	sugarLog = log.Sugar()
}

//func SetGlobalLog(config *LogConfig) {
//	log = NewLog(config)
//	sugarLog = log.Sugar()
//}

func NewLog(config *LogConfig) *zap.Logger {
	if !isDir(config.dir) {
		mkdir(config.dir)
	}
	//todo 不同日志级别下的日志文件不同
	filePath := config.dir + "/apps/" + config.appName + "/" + Level2FileName(config.level) + ".log"

	// lumberjack 实现了一些日志分割的功能
	jLoger := &lumberjack.Logger{
		Filename: filePath,
		MaxSize:  config.maxSize,
	}
	jLoger.LocalTime = config.localTime

	// writer
	var writer zapcore.WriteSyncer
	if config.multiWrite {
		w1 := zapcore.AddSync(jLoger)
		// 日志同时在终端输出
		w2, closeOut, err := zap.Open([]string{"stderr"}...)
		if err != nil {
			if closeOut != nil {
				closeOut()
			}
			panic(err)
		}
		writer = zapcore.NewMultiWriteSyncer(w1, w2)
	} else {
		writer = zapcore.AddSync(jLoger)
	}

	// encoder
	encConfig := zap.NewProductionEncoderConfig()
	encConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var enc zapcore.Encoder
	if config.format == FormatEnumConsole {
		enc = zapcore.NewConsoleEncoder(encConfig)
	} else {
		enc = zapcore.NewJSONEncoder(encConfig)
	}

	// core
	core := zapcore.NewCore(
		enc,
		writer,
		config.level,
	)

	log := zap.New(core)
	//info及以下级别简单打印
	if config.level.Enabled(LevelEnumInfo.Level()) {
		log = log.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.WarnLevel))
	} else {
		log = log.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(config.level))
	}
	if len(config.Hooks) > 0 {
		hks := make([]zap.Option, 0, len(config.Hooks))
		for _, v := range config.Hooks {
			hk := zap.Hooks(v.DoHook)
			hks = append(hks, hk)
		}
		log = log.WithOptions(hks...)
	}
	//日志周期切割
	LogRotatePeriod(jLoger, config.rotatePeriodSecond)
	return log
}

func LogRotatePeriod(log *lumberjack.Logger, periodSecond int64) {
	now := time.Now().Unix()
	//剩余周期时间
	tickSecond := periodSecond - now%periodSecond
	//创建一个定时器，它会在最少过去时间段 d 后到期，向其自身的 C 字段发送当时的时间
	tk := time.NewTicker(time.Duration(tickSecond) * time.Second)
	go func() {
		<-tk.C
		err := log.Rotate()
		if err != nil {
			Error("rotate failed", ErrorField(err))
		}
		tk.Stop()
		periodSecondTicker := time.NewTicker(time.Duration(periodSecond) * time.Second)
		for range periodSecondTicker.C {
			err := log.Rotate()
			if err != nil {
				Error("rotate failed", ErrorField(err))
			}
		}
	}()
}

//创建目录
func mkdir(dir string) {
	if !isDir(dir) {
		err := os.Mkdir(dir, 0777)
		if err != nil {
			panic("mkdir logs failed")
		}
	}
}

//日志级别转化不同文件输出
func Level2FileName(level zap.AtomicLevel) string {
	switch level {
	case LevelEnumDebug, LevelEnumInfo:
		return "info"
	case LevelEnumWarn:
		return "warn"
	case LevelEnumError:
		return "error"
	default:
		return "info"
	}
}

// 判断所给路径是否为文件夹
func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
