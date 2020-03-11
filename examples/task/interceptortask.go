package task

import (
	"queue-task/v1/iface"
	"queue-task/v1/job"
	"queue-task/v1/queue"
	coretask "queue-task/v1/task"
	"queue-task/v1/util"
)

// CreateIntJob 带有拦截器的任务
func CreateIntJob() coretask.CreateJobFunc {
	return func() iface.IJob {
		// 初始化一个队列
		q := queue.NewRedisQueue(&queue.RedisQueueOptions{
			Addr: "redis.dev:6379",
			Key:  "distr debug",
		})
		j := job.NewDefaultJob("int job", q, 10)
		j.RegisterHandleFunc(TestPerform)
		// 注册拦截器
		// 拦截器也可以根据需要进行封，eg:
		// func InterceptorBuilder (v ...interface{} job.WorkInterceptor {
		// 	统一再次注册公共的项目拦截器
		// }
		j.WorkInterceptor(job.ChainWorkInterceptor(WorkInt1, WorkInt2))
		j.SendInterceptor(job.ChainSendInterceptor(SendInt1, SendInt2))
		return j
	}
}

// WorkInt1 业务处理拦截器1
func WorkInt1(data []byte, handle iface.JobHandle) {
	util.WriteLog("work interceptor 1: before handle")
	handle(data)
	util.WriteLog("work interceptor 1: after handle")
}

// SendInt1 发送拦截器1
func SendInt1(msg iface.IMessage, handle iface.SendHandle) {
	util.WriteLog("send interceptor 1: before send")
	handle(msg)
	util.WriteLog("send interceptor 1: after send")
}

// WorkInt2 业务处理拦截器2
func WorkInt2(data []byte, handle iface.JobHandle) {
	util.WriteLog("work interceptor 2: before handle")
	handle(data)
	util.WriteLog("work interceptor 2: after handle")
}

// SendInt2 发送拦截器2
func SendInt2(msg iface.IMessage, handle iface.SendHandle) {
	util.WriteLog("send interceptor 2: before send")
	handle(msg)
	util.WriteLog("send interceptor 2: after send")
}
