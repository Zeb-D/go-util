### go-util

本工具不仅仅对转go开发开封使用，也有高阶或性能用法

对每一个技术点封装整理与特性支持，请看`*_test.go`演示；


### 功能特性

目前对应package说明如下：

> - apollo --读取配置中心
> - bytepool --[]byte 内存池
> - common --通用化工具
> - crypto --加解密封装
> - jsonpath --  [jsonpath README](jsonpath/README.md)
> - cache --  [cache package which implements a fixed-size thread safe LRU cache.](cache/README.md)
> - taskpool -- [goRoutine pool provider core goroutine and keepAliveTime other goroutine](./taskpool/README.md)
> - http --对http相关操作
> - interact --golang实操入门
> - log --封装了zap.Logger 日志操作，支持日志切割
> - mq --简化消息队列中间件操作
> - struct -- [data struct,others contain  algorithm base on it](./struct/README.md)
> - testdata --测试数据

欢迎大家使用，本人增加操作也会向前兼容，请放心使用

### 安装或引入

```
$ go get github.com/Zeb-D/go-util
```

go.mod
```
$ github.com/Zeb-D/go-util v1.0.5
```