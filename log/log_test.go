package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	stdLog "log"
	"os"
	"testing"
)

type hookTest struct{}

func (hookTest) DoHook(e zapcore.Entry) error {
	stdLog.Printf("hook e:%v\n", e)
	return nil
}

func TestSayHello(t *testing.T) {
	println("Hello")
	sayHello("test", "yda")
}

func TestZapInfoLog(t *testing.T) {
	SetGlobalLog("util", false, WithHooksOption(hookTest{}))
	println("zapInfoLog")
	zapInfoLog("info", zap.String("version", "v1"))
	Info("info-test", String("myName", "go-util"))
	Error("error-test", ErrorField(nil))
}

func BenchmarkString(b *testing.B) {
	SetGlobalLog("test", true)
	for i := 0; i < b.N; i++ {
		Info("msg", String("key", "value"))
	}
}

func BenchmarkAny(b *testing.B) {
	SetGlobalLog("test", true)
	for i := 0; i < b.N; i++ {
		Info("msg", Any("key", "value"))
	}
}

func BenchmarkSugar(b *testing.B) {
	SetGlobalLog("test", true)
	for i := 0; i < b.N; i++ {
		sugarLog.Infof("msg key:%s", "value")
	}
}

func TestMKDir(t *testing.T) {
	dir, _ := os.Getwd()
	println("dir:" + dir)
	var testFile = dir + "/test"
	defer os.Remove(testFile)
	mkdir(testFile)
}
