package job

import (
	"fmt"
	"queue-task/v1/conf"
	"queue-task/v1/iface"
	"queue-task/v1/task"
	"queue-task/v1/util"
	"time"
)

// BeforeSendFunc 发送消息前的处理函数
type BeforeSendFunc func()

// AfterSendFunc 发送消息后的处理函数
type AfterSendFunc func()

// BeforeConsumeFunc 消费数据前的处理函数
type BeforeConsumeFunc func()

// AfterConsumeFunc 消费数据后的处理函数
type AfterConsumeFunc func()

// DefaultJob 默认任务
type DefaultJob struct {
	*BaseJob
	workersCnt            int
	exit                  []chan bool
	IsWorking             bool
	beforeSendCallback    BeforeSendFunc
	afterSendCallback     AfterSendFunc
	beforeConsumeCallback BeforeConsumeFunc
	afterConsumeCallback  AfterConsumeFunc
}

// NewDefaultJob default任务构造器
func NewDefaultJob(name string, queue iface.IQueue, conf *conf.DefaultJobConf) *DefaultJob {
	job := &DefaultJob{
		BaseJob:    NewBaseJob(name, queue),
		workersCnt: conf.WorkersCnt,
		exit:       make([]chan bool, conf.WorkersCnt),
	}
	return job
}

// Send 发送消息
func (job *DefaultJob) Send(msg iface.IMessage) {
	if job.beforeSendCallback != nil {
		job.beforeSendCallback()
	}
	ok := job.queue.Enqueue(msg)
	if !ok {
		// write log
		util.WriteLog("send occurred error")
	}
	if job.afterSendCallback != nil {
		job.afterSendCallback()
	}
}

// Work 消费消息
func (job *DefaultJob) Work() {
	if !job.IsWorking {
		task.AddJob(job.name, job)
		job.IsWorking = true
		// 新建channel 因为job是关闭状态
		for i := 0; i < job.workersCnt; i++ {
			job.exit[i] = make(chan bool)
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
	for i := 0; i < job.workersCnt; i++ {
		job.exit[i] <- true
		// 关闭通知chan
		close(job.exit[i])
		util.WriteLog(fmt.Sprintln("default job worker", i, "is stopped"))
	}
	job.IsWorking = false
	job.BaseJob.Stop()
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
		case <-job.exit[id]:
			return
		default:
			if job.beforeConsumeCallback != nil {
				job.beforeConsumeCallback()
			}
			data, ok := job.queue.Dequeue()
			if job.afterConsumeCallback != nil {
				job.afterConsumeCallback()
			}
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

// RegisterBeforeSendFunc 注册发送消息前处理函数
func (job *DefaultJob) RegisterBeforeSendFunc(f BeforeSendFunc) {
	job.beforeSendCallback = f
}

// RegisterAfterSendFunc 注册发送消息后处理函数
func (job *DefaultJob) RegisterAfterSendFunc(f AfterSendFunc) {
	job.afterSendCallback = f
}

// RegisterBeforeConsumeFunc 注册任务处理前回调函数
func (job *DefaultJob) RegisterBeforeConsumeFunc(f BeforeConsumeFunc) {
	job.beforeConsumeCallback = f
}

// RegisterAfterConsumeFunc 注册任务处理后回调函数
func (job *DefaultJob) RegisterAfterConsumeFunc(f AfterConsumeFunc) {
	job.afterConsumeCallback = f
}
