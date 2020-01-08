// Package job 任务接口具体实现包，该包中的job结构实现了IJob接口
//
// 如果需要扩展可将新的job结构代码放在此包中
package job

import (
	"queue-task/v1/iface"
)

// BaseJob 基础job
type BaseJob struct {
	name       string
	queue      iface.IQueue
	handleFunc iface.JobHandle
}

// NewBaseJob 基础job构造器
func NewBaseJob(name string, queue iface.IQueue) *BaseJob {
	return &BaseJob{
		name:  name,
		queue: queue,
	}
}

// Send 发送消息
func (job *BaseJob) Send(msg iface.IMessage) {

}

// Work 分配任务
func (job *BaseJob) Work() {
}

// GetJobName 获得任务名
func (job *BaseJob) GetJobName() string {
	return job.name
}

// GetQueue 获取任务队列
func (job *BaseJob) GetQueue() iface.IQueue {
	return job.queue
}

// IsWorking 是否需要停止
func (job *BaseJob) IsWorking() bool {
	return false
}

// Stop 停止job
func (job *BaseJob) Stop() {
}

// RegisterHandleFunc 注册业务处理方法
func (job *BaseJob) RegisterHandleFunc(f iface.JobHandle) {
	job.handleFunc = f
}
