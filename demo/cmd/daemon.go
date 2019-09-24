package main

import (
	"path/filepath"
	"queue-task/demo/task"
	"queue-task/v1/conf"
	coretask "queue-task/v1/task"
	"queue-task/v1/util"
	"sync"
	"time"
)

var once sync.Once

func main() {
	confPath, _ := filepath.Abs("../conf")
	conf.Init(confPath)
	handleCreateFunc()
	// 测试任务启动，运行
	coretask.Work()
	// for _, function := range util.CreateFuncList {
	// 	job := function()
	// 	job.Work()
	// }
	// 测试消息发送
	task.TestSendMsg()
	time.Sleep(5 * time.Second)
	for _, job := range util.JobList {
		job.Stop()
	}
}

func handleCreateFunc() {
	once.Do(func() {
		coretask.AddCreateFunc("test", task.CreateTestJob())
	})
}
