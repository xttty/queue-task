package queue

import (
	"fmt"
	"queue-task/v1/iface"
	"queue-task/v1/util"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// KafkaQueue kafka队列
type KafkaQueue struct {
	server  string
	group   string
	topic   string
	timeout time.Duration
}

// NewKafkaQueue 新建kafka队列实例
func NewKafkaQueue(server, group, topic string, timeout time.Duration) (*KafkaQueue, error) {
	kq := &KafkaQueue{}
	kq.server = server
	kq.group = group
	kq.topic = topic
	if timeout <= 0 {
		timeout = -1
	}
	kq.timeout = timeout
	return kq, nil
}

// Enqueue 写入消息
func (kq *KafkaQueue) Enqueue(msg iface.IMessage) bool {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kq.server,
	})
	if err != nil {
		util.WriteLog(fmt.Sprintln("new kafka produce failed, err:", err.Error()))
		return false
	}
	defer p.Close()

	data, err := msg.ToJSON()
	if err != nil {
		util.WriteLog(fmt.Sprintln("message to json failed, err:", err.Error()))
		return false
	}
	pChan := make(chan kafka.Event)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kq.topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}, pChan)
	// 生产结果消息返回
	select {
	case e := <-pChan:
		deliverMsg := e.(*kafka.Message)
		if deliverMsg.TopicPartition.Error != nil {
			util.WriteLog(fmt.Sprintln("kafka produce failed, err:", deliverMsg.TopicPartition.Error.Error()))
			return false
		}
	case <-time.Tick(kq.timeout):
		util.WriteLog(fmt.Sprintln("kafka produce time-out"))
		return false
	}
	close(pChan)
	return true
}

// Dequeue 取出消息
func (kq *KafkaQueue) Dequeue() ([]byte, bool) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kq.server,
		"group.id":          kq.group,
		"auto.offset.reset": "beginning",
	})
	if err != nil {
		util.WriteLog(fmt.Sprintln("new kafka consumer failed, err:", err.Error()))
		return nil, false
	}
	defer c.Close()

	topics := []string{kq.topic}
	c.SubscribeTopics(topics, nil)
	msg, err := c.ReadMessage(kq.timeout)
	if err != nil {
		util.WriteLog(fmt.Sprintln("kafka consume failed, err", err.Error()))
		return nil, false
	}
	return msg.Value, true
}

// Size 消息lag量
func (kq *KafkaQueue) Size() int64 {
	return int64(-1)
}

func (kq *KafkaQueue) Debug() {
	c, _ := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":      kq.server,
		"group.id":               kq.group,
		"auto.offset.reset":      "earliest",
		"statistics.interval.ms": 5000,
	})
	defer c.Close()
	// low, high, err := c.QueryWatermarkOffsets(kq.topic, 0, int(time.Second))
	// util.WriteLog(fmt.Sprintf("low: %d, high: %d, err: %s", low, high, err.Error()))
	committedOffsets, err := c.Committed([]kafka.TopicPartition{{
		Topic:     &kq.topic,
		Partition: int32(0),
	}}, int(3*time.Second))
	if err != nil {
		util.WriteLog(fmt.Sprintln(err.Error()))
	}
	for _, info := range committedOffsets {
		util.WriteLog(fmt.Sprintln("partition:", info.Partition, "offset", info.Offset, "metaData:", info.Metadata))
	}
}
