package queue_test

import (
	"encoding/json"
	"queue-task/v1/msg"
	"queue-task/v1/queue"
	"queue-task/v1/util"
	"testing"
	"time"
)

func TestKafka(t *testing.T) {
	util.RegisterLogHandle(func(str string) {
		t.Log(str)
	})
	// 不知道为何消费kafka需要4秒，不然就会timeout
	kafkaQ, err := queue.NewKafkaQueue("dockervm:9092", "test-consumer-group", "xty", 4*time.Second, 2)
	if err != nil {
		t.Log(err)
		return
	}
	message := &msg.BaseMsg{
		Data: msg.H{
			"string": "hello queue task",
		},
	}
	ok := kafkaQ.Enqueue(message)
	t.Log("enqueue result", ok)
	time.Sleep(time.Second)
	data, ok := kafkaQ.Dequeue()
	deMsg := msg.BaseMsg{}
	json.Unmarshal(data, &deMsg)
	t.Log("dequeue result:", deMsg, ok)
}
