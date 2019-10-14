package job

import (
	"context"
	"queue-task/v1/iface"
	"time"
)

// DistributedJob 分布式任务
// 在分布式环境下依然能保持正常并发的任务执行
type DistributedJob struct {
	baseJob     *BaseJob
	rootCancel  context.CancelFunc
	childCancel map[int]context.CancelFunc
	WorkerDur   time.Duration //worker协程每次工作的时长
	queue       iface.IQueue
}
