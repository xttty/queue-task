package job

import (
	"context"
	"fmt"
	"queue-task/v1/conf"
	"queue-task/v1/iface"
	"queue-task/v1/task"
	"queue-task/v1/util"
	"time"
)

// DefaultJob 默认任务
type DefaultJob struct {
	*BaseJob
	workersCnt     int
	IsWorking      bool
	ctx            context.Context
	ctxCancel      context.CancelFunc
	childCtx       []context.Context
	childCtxCancel []context.CancelFunc
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
	if !job.IsWorking {
		job.workPrepare()
		job.IsWorking = true
		for i := 0; i < job.workersCnt; i++ {
			job.childCtx[i], job.childCtxCancel[i] = context.WithCancel(job.ctx)
		}
		for i := 0; i < job.workersCnt; i++ {
			go job.startWorker(i, job.BaseJob.handleFunc)
		}
	}
}

// Stop 停止消费
func (job *DefaultJob) Stop() {
	if !job.IsWorking {
		return
	}
	job.ctxCancel()
	job.IsWorking = false
	job.BaseJob.Stop()
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

func (job *DefaultJob) startWorker(id int, f iface.JobHandle) {
	util.WriteLog(fmt.Sprintln("default job worker", id, "is starting"))
	var timeRest time.Duration
	for {
		select {
		case <-job.childCtx[id].Done():
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

func (job *DefaultJob) workPrepare() {
	task.AddJob(job.name, job)
	job.ctx, job.ctxCancel = context.WithCancel(context.Background())
	job.childCtx = make([]context.Context, job.workersCnt)
	job.childCtxCancel = make([]context.CancelFunc, job.workersCnt)
}
