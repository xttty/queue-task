package main

import (
	"fmt"
	"queue-task/examples/task"
	coretask "queue-task/v1/task"
	"queue-task/v1/util"
	"sync"
	"time"
)

var once sync.Once

func main() {
	handleCreateFunc()
	util.RegisterLogHandle(func(s string) {
		fmt.Println(s)
	})
	// 测试任务启动，运行
	coretask.Work()
	// 测试消息发送
	task.TestSendMsg()
	time.Sleep(5 * time.Second)
	coretask.Stop()
	time.Sleep(time.Second)
}

func handleCreateFunc() {
	once.Do(func() {
		coretask.AddCreateFunc("test", task.CreateTestJob())
	})
}
