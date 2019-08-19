// Package task 任务调度包，所有任务都在该包中
package task

import "queue-task/v1/iface"

// CreateJobFunc 新建job方法
type CreateJobFunc func() iface.IJob

// JobList 任务表
var JobList = map[string]iface.IJob{}

// CreateFuncList 注册job创建方法
var CreateFuncList = map[string]CreateJobFunc{
	"test": CreateTestJob(),
}
