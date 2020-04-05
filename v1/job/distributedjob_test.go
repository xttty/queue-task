package job

import (
	"encoding/json"
	"math/rand"
	"queue-task/v1/queue"
	"queue-task/v1/util"
	"testing"
	"time"
)

type TestMsg struct {
	Key  string
	Time string
}

func (msg *TestMsg) Serialize() ([]byte, error) {
	data, err := json.Marshal(*msg)
	return data, err
}

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
		message := &TestMsg{}
		json.Unmarshal(data, message)
		t.Log(*message)
	}
	job := NewDistributedJob("debug", q, SetCircleTime(5*time.Second), SetWorkerStrategy(cntFunc))
	job.RegisterHandleFunc(handleFunc)
	job.Work()
	for i := 0; i < 10; i++ {
		job.Send(&TestMsg{
			Key:  "debug",
			Time: time.Now().String(),
		})
	}
	time.Sleep(30 * time.Second)
	job.Stop()
	job.Work()
	time.Sleep(2 * time.Second)
	job.Stop()
	time.Sleep(time.Second)
}
