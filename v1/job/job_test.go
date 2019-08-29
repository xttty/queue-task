package job_test

import (
	"queue-task/v1/conf"
	"queue-task/v1/job"
	"queue-task/v1/queue"
	"testing"
	"time"
)

func TestJob(t *testing.T) {
	qConf := conf.RedisQueueConf{
		Addr: "redis.dev:6379",
	}
	jConf := conf.DefaultJobConf{
		WorkersCnt: 5,
	}
	q := queue.NewRedisQueue(&qConf, "test")
	testJob := job.NewDefaultJob("test", q, &jConf)
	testJob.Work()
	time.Sleep(3 * time.Second)
	testJob.Stop()
}
