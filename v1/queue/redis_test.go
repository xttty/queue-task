package queue_test

import (
	"encoding/json"
	"queue-task/v1/msg"
	"queue-task/v1/queue"
	"testing"
)

func TestTest(t *testing.T) {
	testMsg := &msg.BaseMsg{
		Data: msg.H{
			"field1": "data1",
			"field2": "data2",
		},
	}
	redis := queue.NewRedisQueue(&queue.RedisQueueOptions{
		Addr: "redis.dev:6379",
		Key:  "debug",
	})
	if ok := redis.Enqueue(testMsg); ok {
		t.Log("send msg succes")
	} else {
		t.Log("send msg failed")
	}

	len := redis.Size()
	t.Log(len)

	if data, ok := redis.Dequeue(); ok {
		message := msg.BaseMsg{}
		json.Unmarshal(data, &message)
		t.Log(message)
	} else {
		t.Log("get msg failed")
	}

}
