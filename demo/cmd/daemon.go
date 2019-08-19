package main

import (
	"fmt"
	"path/filepath"
	"queue-task/demo/conf"
	"queue-task/demo/task"
	"queue-task/v1/msg"
	"time"
)

func main() {
	rootPath, _ := filepath.Abs("../../")
	conf.Init(rootPath)
	for name, function := range task.CreateFuncList {
		job := function()
		task.JobList[name] = job
		job.Work()
	}
	// 测试job
	// test()
	time.Sleep(2 * time.Second)
	for _, job := range task.JobList {
		job.Stop()
	}
}

func test() {
	job := task.JobList["test"]
	fmt.Println(job.GetJobName())
	for i := 0; i < 10; i++ {
		msg := &msg.BaseMsg{
			Data: msg.H{
				"id": i,
			},
		}
		job.Send(msg)
	}
}
