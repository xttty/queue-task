package job

import (
	"fmt"
	"queue-task/v1/conf"
	"queue-task/v1/iface"
	"queue-task/v1/util"
	"time"
)

// DefaultJob kafka任务
type DefaultJob struct {
	*BaseJob
	workersCnt int
	exit       []chan bool
	IsWorking  bool
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
	ok := job.queue.Enqueue(msg)
	if !ok {
		// write log
		util.WriteLog("send occurred error")
	}
}

// Work 消费消息
func (job *DefaultJob) Work() {
	if !job.IsWorking {
		job.BaseJob.Work()
		job.IsWorking = true
		// 新建channel 因为job是关闭状态
		for i := 0; i < job.workersCnt; i++ {
			job.exit[i] = make(chan bool, 1)
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
		fmt.Println("kafka job worker", i, "is stopped")
	}
	job.IsWorking = false
	job.BaseJob.Stop()
}

// RegisterHandleFunc 注册业务回调方法
func (job *DefaultJob) RegisterHandleFunc(f iface.JobHandle) {
	job.BaseJob.RegisterHandleFunc(f)
}

func (job *DefaultJob) startWorker(id int, f iface.JobHandle) {
	fmt.Println("kafka job worker", id, "is starting")
	var timeRest time.Duration
	for {
		select {
		case <-job.exit[id]:
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
