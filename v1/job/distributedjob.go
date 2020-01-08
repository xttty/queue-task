package job

import (
	"queue-task/v1/iface"
	"sync"
	"time"
)

// CurrentWorkerCntFunc 获取当前worker数方法
type CurrentWorkerCntFunc func() int

// DistributedJobOption  分布式任务设置方法
type DistributedJobOption func(*DistributedJob) *DistributedJob

// DistributedJob 分布式任务
// 在分布式环境下依然能保持正常并发的任务执行
type DistributedJob struct {
	*DefaultJob
	circleTime   time.Duration
	getWorkerCnt CurrentWorkerCntFunc
	exitCh       chan struct{}
	rwLock       sync.RWMutex
	isWorking    bool
}

var defaultCircleTime = 60 * time.Second
var defaultGetWorkerCnt = func() int {
	return 0
}

// NewDistributeJob 新建一个分布式任务
func NewDistributeJob(name string, q iface.IQueue, options ...DistributedJobOption) *DistributedJob {
	job := &DistributedJob{
		DefaultJob:   NewDefaultJob(name, q, 0),
		circleTime:   defaultCircleTime,
		getWorkerCnt: defaultGetWorkerCnt,
		exitCh:       make(chan struct{}, 0),
		isWorking:    false,
	}
	for i := 0; i < len(options); i++ {
		job = options[i](job)
	}
	return job
}

// SetCircleTime 设置轮询时间
func SetCircleTime(t time.Duration) DistributedJobOption {
	return func(job *DistributedJob) *DistributedJob {
		if t >= time.Second {
			job.circleTime = t
		}
		return job
	}
}

// SetWorkCntFunc 设置并发控制函数
func SetWorkCntFunc(f CurrentWorkerCntFunc) DistributedJobOption {
	return func(job *DistributedJob) *DistributedJob {
		if f != nil {
			job.getWorkerCnt = f
		}
		return job
	}
}

// Work 分布式任务工作方法
func (job *DistributedJob) Work() {
	if !job.IsWorking() {
		job.rwLock.Lock()
		job.exitCh = make(chan struct{})
		job.isWorking = true
		job.rwLock.Unlock()
	} else {
		return
	}
	go func() {
		ch := job.exitCh
		if ch == nil {
			return
		}
		for {
			select {
			case <-time.Tick(job.circleTime):
				cnt := job.getWorkerCnt()
				if cnt < 0 {
					cnt = 0
				}
				if cnt != job.DefaultJob.workersCnt {
					job.DefaultJob.ChangeWorkerCnt(cnt)
				}
			case <-ch:
				return
			}
		}
	}()
}

// Stop 分布式任务停止方法
func (job *DistributedJob) Stop() {
	if !job.IsWorking() {
		return
	}
	job.rwLock.Lock()
	close(job.exitCh)
	job.exitCh = nil
	job.isWorking = false
	job.rwLock.Unlock()
	job.DefaultJob.Stop()
}

// IsWorking 判断任务是否正在运行
func (job *DistributedJob) IsWorking() bool {
	job.rwLock.RLock()
	defer job.rwLock.RUnlock()

	return job.isWorking
}
