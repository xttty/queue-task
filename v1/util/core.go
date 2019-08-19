package util

import (
	"fmt"
	"sync"
)

// LogHandleFunc 日志处理函数
type LogHandleFunc func(string)

// Core ...
type Core struct {
	logFunc LogHandleFunc
}

var glCore *Core
var once sync.Once

// Instance 核心配置单例
func Instance() *Core {
	once.Do(func() {
		glCore = &Core{
			logFunc: defaultLogFunc,
		}
	})
	return glCore
}

// RegisterLogHandle 注册日志记录方法
func RegisterLogHandle(f LogHandleFunc) {
	Instance().logFunc = f
}

// WriteLog 写入日志
func WriteLog(s string) {
	Instance().logFunc(s)
}

func defaultLogFunc(s string) {
	fmt.Println(s)
}
