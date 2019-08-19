package main

import (
	"path/filepath"
	"queue-task/demo/conf"
	"queue-task/demo/task"
	"queue-task/v1/util"
	"sync"
	"time"
)

var once sync.Once

func main() {
	confPath, _ := filepath.Abs("../conf")
	conf.Init(confPath)
	handleCreateFunc()
	for _, function := range util.CreateFuncList {
		job := function()
		job.Work()
	}
	// 测试job
	time.Sleep(2 * time.Second)
	for _, job := range util.JobList {
		job.Stop()
	}
}

func handleCreateFunc() {
	once.Do(func() {
		util.AddCreateFunc("test", task.CreateTestJob())
	})
}
