package iface

// IQueue 队列接口
type IQueue interface {
	Enqueue(msg IMessage) bool // 写入队列
	Dequeue() ([]byte, bool)   // 读取队列
	Size() int64               // 队列长度
}
