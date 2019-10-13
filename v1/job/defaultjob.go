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
	workMidwares   []WorkMidwareFunc
	sendMidwares   []SendMidwareFunc
	rwLock         *sync.RWMutex
}

// WorkMidwareFunc job运行中间件
// 中间件一般用于数据统计、日志等
// 每个work中间件都必须显示调用 job.WorkNext(data, id+1) 不然中间件无法正常执行
type WorkMidwareFunc func(job *DefaultJob, data []byte, id int)

// SendMidwareFunc job send中间件
// 每个send中间件都必须显示调用 job.SendNext(msg, id+1) 不然中间件无法正常执行
type SendMidwareFunc func(job *DefaultJob, msg iface.IMessage, id int)

// NewDefaultJob default任务构造器
func NewDefaultJob(name string, queue iface.IQueue, conf *conf.DefaultJobConf) *DefaultJob {
	job := &DefaultJob{
		BaseJob:    NewBaseJob(name, queue),
		workersCnt: conf.WorkersCnt,
		rwLock:     new(sync.RWMutex),
		IsWorking:  false,
	}
	job.workMidwares = append(job.workMidwares, doWork)
	job.sendMidwares = append(job.sendMidwares, doSend)
	return job
}

// Send 发送消息
func (job *DefaultJob) Send(msg iface.IMessage) {
	job.SendNext(msg, 0)
}

// SendNext send中间件控制，使用方式和work中间件一样
func (job *DefaultJob) SendNext(msg iface.IMessage, idx int) {
	if idx < len(job.sendMidwares) {
		job.sendMidwares[idx](job, msg, idx)
	}
}

func doSend(job *DefaultJob, msg iface.IMessage, idx int) {
	ok := job.queue.Enqueue(msg)
	if !ok {
		// write log
		util.WriteLog("send occurred error")
	}
	job.SendNext(msg, idx+1)
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
			go job.startWorker(childCtx, i)
		}
		job.IsWorking = true
		job.rwLock.Unlock()
		return
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
	job.rwLock.RUnlock()
	job.rwLock.Lock()
	// 通知子worker结束
	job.ctxCancel()
	// 基础任务也执行stop
	job.BaseJob.Stop()
	// 标记运行状态为结束
	job.IsWorking = false
	util.WriteLog(fmt.Sprintf("job[%s] is stopped", job.GetJobName()))
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
				// 开始执行work中间件
				job.WorkNext(data, 0)
				timeRest = 100 * time.Microsecond
			} else {
				timeRest = 1000 * time.Microsecond
			}
		}
		time.Sleep(timeRest)
	}
}

func doWork(job *DefaultJob, data []byte, id int) {
	job.BaseJob.handleFunc(data)
}

// WorkNext 中间件控制
// func midware(job *DefaultJob, data []byte, idx int) {
// 	提前于handle方法的逻辑
// 	......
// 	job.WorkNext(data, idx+1) 必须要有这个方法
// 	 handle方法处理完毕后的逻辑
// 	......
// }
func (job *DefaultJob) WorkNext(data []byte, idx int) {
	if idx < len(job.workMidwares) {
		job.workMidwares[idx](job, data, idx)
	}
}

// RegisterWorkMidware 注册work中间件
func (job *DefaultJob) RegisterWorkMidware(midwares ...WorkMidwareFunc) {
	job.workMidwares = append(midwares, job.workMidwares...)
}

// RegisterSendMidware 注册send中间件
func (job *DefaultJob) RegisterSendMidware(midwares ...SendMidwareFunc) {
	job.sendMidwares = append(midwares, job.sendMidwares...)
}
