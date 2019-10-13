package job

import (
	"context"
	"fmt"
	"queue-task/v1/conf"
	"queue-task/v1/iface"
	"queue-task/v1/util"
	"sync"
	"time"
)

// DefaultJob 默认任务
// 保留worker的子context，方便以后单独管理单个worker
type DefaultJob struct {
	*BaseJob
	workersCnt     int
	IsWorking      bool
	ctxCancel      context.CancelFunc
	childCtxCancel []context.CancelFunc
	rwLock         sync.RWMutex
}

// NewDefaultJob default任务构造器
func NewDefaultJob(name string, queue iface.IQueue, conf *conf.DefaultJobConf) *DefaultJob {
	job := &DefaultJob{
		BaseJob:    NewBaseJob(name, queue),
		workersCnt: conf.WorkersCnt,
	}
	return job
}

// Send 发送消息
func (job *DefaultJob) Send(msg iface.IMessage) {
	ok := job.queue.Enqueue(msg)
	if !ok {
		// write log
		util.WriteLog("send occurred error")
	}
}

// Work 消费消息
func (job *DefaultJob) Work() {
	// 并发控制
	job.rwLock.RLock()
	if !job.IsWorking {
		job.rwLock.RUnlock()
		job.rwLock.Lock()
		var jobCtx context.Context
		jobCtx, job.ctxCancel = context.WithCancel(context.Background())
		job.childCtxCancel = make([]context.CancelFunc, job.workersCnt)
		for i := 0; i < job.workersCnt; i++ {
			var childCtx context.Context
			childCtx, job.childCtxCancel[i] = context.WithCancel(jobCtx)
			go job.startWorker(childCtx, i, job.BaseJob.handleFunc)
		}
		job.IsWorking = true
		job.rwLock.Unlock()
	}
	job.rwLock.RUnlock()
}

// Stop 停止消费
func (job *DefaultJob) Stop() {
	// 并发控制
	job.rwLock.RLock()
	if !job.IsWorking {
		job.rwLock.RUnlock()
		return
	}
	job.rwLock.Unlock()
	job.rwLock.Lock()
	job.ctxCancel()
	job.IsWorking = false
	job.BaseJob.Stop()
	job.rwLock.Unlock()
}

// ChangeWorkerCnt 改变并发数，并重启
func (job *DefaultJob) ChangeWorkerCnt(cnt int) {
	if job.workersCnt != cnt {
		job.Stop()
		job.workersCnt = cnt
		job.Work()
	}
}

// GetWorkerCnt 查看并发数
func (job *DefaultJob) GetWorkerCnt() int {
	return job.workersCnt
}

// RegisterHandleFunc 注册业务回调方法
func (job *DefaultJob) RegisterHandleFunc(f iface.JobHandle) {
	job.BaseJob.RegisterHandleFunc(f)
}

func (job *DefaultJob) startWorker(ctx context.Context, id int, f iface.JobHandle) {
	util.WriteLog(fmt.Sprintln("default job worker", id, "is starting"))
	var timeRest time.Duration
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, ok := job.queue.Dequeue()
			if ok {
				f(data)
				timeRest = 100 * time.Microsecond
			} else {
				timeRest = 1000 * time.Microsecond
			}
		}
		time.Sleep(timeRest)
	}
}
