package task

import (
	"queue-task/v1/iface"
	"queue-task/v1/job"
	"queue-task/v1/queue"
	coretask "queue-task/v1/task"
	"time"
)

// CreateDistrJob 分布式任务构建方法
func CreateDistrJob() coretask.CreateJobFunc {
	return func() iface.IJob {
		// 初始化一个队列
		q := queue.NewRedisQueue(&queue.RedisQueueOptions{
			Addr: "redis.dev:6379",
			Key:  "distr debug",
		})
		// 生成job
		// 控制并发的策略
		// 在实际项目中该策略应该独立出来，eg:
		// func getStrategy(name string) job.WorkerStrategyFunc {
		// 	 ...
		// 	return  func () int {

		// 	}
		// }
		strategy := func() int {
			return 5
		}
		j := job.NewDistributedJob("distr job", q, job.SetCircleTime(10*time.Second), job.SetWorkerStrategy(strategy))
		j.RegisterHandleFunc(TestPerform)
		return j
	}
}
