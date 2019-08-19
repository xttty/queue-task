// Package queue 队列具体实现包，实现iface.IQueue接口
// 更多的队列实现都放在改包中，前提是必须实现iface.IQueue接口
package queue

import (
	"fmt"
	"queue-task/conf"
	"queue-task/v1/iface"

	"gitlab.meitu.com/platform/kaproxy/client"
)

// Kafka kafka队列
type Kafka struct {
	client *client.KaproxyClient
	host   string
	port   int
	token  string
	topic  string
	group  string
}

// NewKaproxy 生成Kaproxy实例
func NewKaproxy(conf conf.KaproxyQueueConf, topic, group string) *Kafka {
	return &Kafka{
		client: client.NewKaproxyClient(conf.Host, conf.Port, conf.Token),
		topic:  topic,
		group:  group,
	}
}

// Enqueue kafka写入
func (q *Kafka) Enqueue(msg iface.IMessage) bool {
	data, jsonErr := msg.ToJSON()
	if jsonErr != nil {
		// writelog
		fmt.Println("msg to json occurred error", jsonErr)
		return false
	}
	kafkaMsg := client.Message{
		Value: data,
	}
	resp, prodErr := q.client.Produce(q.topic, kafkaMsg)
	if prodErr != nil {
		fmt.Println("kafka queue produce occurred error:", prodErr)
		return false
	}
	if resp != nil {
		fmt.Println("kafka queue produce response:", resp)
	}
	return true
}

// ReverseEnqueue kafka逆序写入
func (q *Kafka) ReverseEnqueue(msg iface.IMessage) bool {
	return q.Enqueue(msg)
}

// Dequeue kafka队列dequeue方法
func (q *Kafka) Dequeue() ([]byte, bool) {
	resp, err := q.client.Consume(q.group, q.topic)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return resp.Value, true
}

// Size 队列长度
func (q *Kafka) Size() int {
	return 0
}
