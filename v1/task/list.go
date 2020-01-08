package task

import (
	"queue-task/v1/iface"
	"time"
)

// CreateJobFunc 新建job方法
type CreateJobFunc func() iface.IJob

// jobList 任务表
var jobList = map[string]iface.IJob{}

// CreateFuncList job创建方法列表
var CreateFuncList = map[string]CreateJobFunc{}

// AddCreateFunc 注册job创建方法
func AddCreateFunc(key string, f CreateJobFunc) {
	CreateFuncList[key] = f
}

// DelCreateFunc 删除job创建方法
func DelCreateFunc(key string) {
	delete(CreateFuncList, key)
}

// AddJob 添加任务
func AddJob(key string, job iface.IJob) {
	jobList[key] = job
}

// DelJob 将job从list中删除
func DelJob(key string) {
	delete(jobList, key)
}

// Work 任务启动
func Work() {
	for _, createFunc := range CreateFuncList {
		job := createFunc()
		AddJob(job.GetJobName(), job)
		job.Work()
	}
}

// Stop 任务停止
func Stop() {
	for _, job := range jobList {
		job.Stop()
		DelJob(job.GetJobName())
	}
	// TODO 有必要使用sync.waitGroup或者exit channel 来控制退出吗？
	// 等待1S让任务协程有退出的时间
	time.Sleep(1 * time.Second)
}

// Restart 重启
func Restart() {
	Stop()
	Work()
}
