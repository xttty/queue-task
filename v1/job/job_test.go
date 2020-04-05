package job

import (
	"queue-task/v1/iface"
	"queue-task/v1/msg"
	"queue-task/v1/queue"
	"testing"
	"time"
)

func TestJob(t *testing.T) {
	q := queue.NewRedisQueue(&queue.RedisQueueOptions{
		Addr: "redis.dev:6379",
		Key:  "debug",
	})

	print1WorkInterceptor := func(value []byte, handle iface.JobHandle) {
		t.Log("print1 before work handle")
		handle(value)
		t.Log("print1 after work handle")
	}

	print2WorkInterceptor := func(value []byte, handle iface.JobHandle) {
		t.Log("print2 before work handle")
		handle(value)
		t.Log("print2 after work handle")
	}

	print3SendInterceptor := func(msg iface.IMessage, handle iface.SendHandle) {
		t.Log("print 3 before send handle")
		handle(msg)
		t.Log("print 3 after send handle")
	}

	print4SendInterceptor := func(msg iface.IMessage, handle iface.SendHandle) {
		t.Log("print 4 before send handle")
		handle(msg)
		t.Log("print 4 after send handle")
	}

	testJob := NewDefaultJob("test", q, 10)
	testJob.WorkInterceptor(ChainWorkInterceptor(print1WorkInterceptor, print2WorkInterceptor))
	testJob.SendInterceptor(ChainSendInterceptor(print3SendInterceptor, print4SendInterceptor))
	testJob.RegisterHandleFunc(func(data []byte) {
		t.Log(data)
	})
	go testJob.Send(&msg.BaseMsg{
		Data: msg.H{
			"time":    123,
			"content": "xty_debug",
		},
	})
	testJob.Work()
	time.Sleep(3 * time.Second)
	testJob.Stop()
	time.Sleep(1 * time.Second)
}
