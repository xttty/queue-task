package job_test

import (
	"math/rand"
	"queue-task/v1/job"
	"queue-task/v1/msg"
	"queue-task/v1/queue"
	"queue-task/v1/util"
	"testing"
	"time"
)

func TestDistributeJob(t *testing.T) {
	util.RegisterLogHandle(func(s string) {
		t.Log(s)
	})
	q := queue.NewRedisQueue(&queue.RedisQueueOptions{
		Addr: "redis.dev:6379",
		Key:  "debug",
	})
	rand.Seed(time.Now().UnixNano())
	cntFunc := func() int {
		cnt := rand.Intn(5) + 1
		t.Log("cnt:", cnt)
		return cnt
	}
	handleFunc := func(data []byte) {
		t.Log(data)
	}
	job := job.NewDistributeJob("debug", q, job.SetCircleTime(5*time.Second), job.SetWorkCntFunc(cntFunc))
	job.RegisterHandleFunc(handleFunc)
	job.Work()
	for i := 0; i < 10; i++ {
		job.Send(&msg.BaseMsg{
			Data: msg.H{
				"test": time.Now(),
			},
		})
	}
	time.Sleep(30 * time.Second)
	job.Stop()
	job.Work()
	time.Sleep(2 * time.Second)
	job.Stop()
	time.Sleep(time.Second)
}
