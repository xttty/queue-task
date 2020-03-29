package job

import (
	"context"
	"fmt"
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
	isWorking      bool
	ctxCancel      context.CancelFunc
	childCtxCancel []context.CancelFunc
	rwLock         *sync.RWMutex
	workInt        WorkInterceptor
	sendInt        SendInterceptor
}

// NewDefaultJob default任务构造器
func NewDefaultJob(name string, queue iface.IQueue, workCnt int) *DefaultJob {
	job := &DefaultJob{
		BaseJob:    NewBaseJob(name, queue),
		workersCnt: workCnt,
		rwLock:     new(sync.RWMutex),
		isWorking:  false,
	}
	return job
}

// Send 发送消息
func (job *DefaultJob) Send(msg iface.IMessage) {
	if job.sendInt == nil {
		job.doSend(msg)
	} else {
		job.sendInt(msg, job.doSend)
	}
}

func (job *DefaultJob) doSend(msg iface.IMessage) {
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
	defer job.rwLock.Unlock()

	if !job.isWorking {
		var jobCtx context.Context
		jobCtx, job.ctxCancel = context.WithCancel(context.Background())
		job.childCtxCancel = make([]context.CancelFunc, job.workersCnt)
		for i := 0; i < job.workersCnt; i++ {
			var childCtx context.Context
			childCtx, job.childCtxCancel[i] = context.WithCancel(jobCtx)
			go job.startWorker(childCtx, i)
		}
		job.isWorking = true
		return
	}
}

// IsWorking 任务是否在运行
func (job *DefaultJob) IsWorking() bool {
	return job.isWorking
}

// Stop 停止消费
func (job *DefaultJob) Stop() {
	// 并发控制
	job.rwLock.Lock()
	defer job.rwLock.Unlock()

	if !job.isWorking {
		return
	}
	// 通知子worker结束
	job.ctxCancel()
	// 基础任务也执行stop
	job.BaseJob.Stop()
	// 标记运行状态为结束
	job.isWorking = false
	util.WriteLog(fmt.Sprintf("job[%s] is stopped", job.GetJobName()))
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

func (job *DefaultJob) startWorker(ctx context.Context, id int) {
	util.WriteLog(fmt.Sprintf("job[%s] worker %d is starting\n", job.GetJobName(), id))
	defer util.WriteLog(fmt.Sprintf("job[%s] worker %d is stopped\n", job.GetJobName(), id))

	var timeRest time.Duration
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, ok := job.queue.Dequeue()
			if ok {
				if job.workInt == nil {
					job.handleFunc(data)
				} else {
					job.workInt(data, job.handleFunc)
				}
				timeRest = 100 * time.Microsecond
			} else {
				timeRest = 1000 * time.Microsecond
			}
		}
		time.Sleep(timeRest)
	}
}

// WorkInterceptor 注册work拦截器
func (job *DefaultJob) WorkInterceptor(workInt WorkInterceptor) {
	job.workInt = workInt
}

// SendInterceptor 注册send拦截器
func (job *DefaultJob) SendInterceptor(sendInt SendInterceptor) {
	job.sendInt = sendInt
}
