// Package iface 接口定义归档
package iface

// IJob 任务接口
type IJob interface {
	Send(msg IMessage)              // 发送消息
	Work()                          // 执行业务
	Stop()                          // 停止任务
	GetJobName() string             // 获取队列名
	NeedStop() bool                 // 是否停止
	GetQueue() IQueue               // 获取队列
	RegisterHandleFunc(f JobHandle) //注册业务处理方法
}

// JobHandle 业务处理方法
type JobHandle func([]byte)

// SendHandle send函数
type SendHandle func(msg IMessage)
