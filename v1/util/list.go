package util

import "queue-task/v1/iface"

// CreateJobFunc 新建job方法
type CreateJobFunc func() iface.IJob

// JobList 任务表
var JobList = map[string]iface.IJob{}

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
	JobList[key] = job
}

// DelJob 将job从list中删除
func DelJob(key string) {
	delete(JobList, key)
}
