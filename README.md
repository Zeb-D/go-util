### go-util

考虑平时对go进行coding，会涉及到一些第三方库集成，争取做到开封即用原则；

对每一个技术点封装整理与特性支持，请看`*_test.go`演示；



### 目录说明

本项目主要是各种工具封装与整理，但目录言简意赅

> - apollo --读取配置中心
> - bytepool --[]byte 内存池
> - common --通用化工具
> - crypto --加解密封装
> - cache --  [cache package which implements a fixed-size thread safe LRU cache.](./cache/README.md)
> - taskpool -- [goRoutine pool provider core goroutine and keepAliveTime other goroutine](./taskpool/README.md)
> - http --对http相关操作
> - interact --自己平时实操golang的痕迹
> - log --封装了zap.Logger 日志操作，支持日志切割
> - mq --简化消息队列中间件操作
> - redis --简化redis封装
> - struct -- [data struct,others contain  algorithm base on it](./struct/README.md)
> - testdata --测试数据
> - todo --在实验操作中，将要整理成目录

欢迎大家使用，本人增加操作也会向前兼容，请放心使用

