package queue

import (
	"fmt"
	"queue-task/v1/iface"
	"queue-task/v1/util"

	"github.com/go-redis/redis"
)

// RedisQueue 基于redis的队列
type RedisQueue struct {
	client *redis.Client
	key    string
}

// NewRedisQueue 构造redis队列
func NewRedisQueue(addr, password, key string, db int) *RedisQueue {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisQueue{
		client: client,
		key:    key,
	}
}

// Enqueue 实现iqueue接口的写入方法
func (rq *RedisQueue) Enqueue(msg iface.IMessage) bool {
	data, err := msg.ToJSON()
	if err != nil {
		util.WriteLog(err.Error())
		return false
	}
	err = rq.client.LPush(rq.key, data).Err()
	if err != nil {
		util.WriteLog(fmt.Sprintf("enqueue failed\tkey:%sdata:%s\nerr:%s", rq.key, data, err.Error()))
		return false
	}
	return true
}

// Dequeue 实现iface接口读取方法
func (rq *RedisQueue) Dequeue() ([]byte, bool) {
	data, err := rq.client.LPop(rq.key).Result()
	if err != nil {
		util.WriteLog(fmt.Sprintf("dequeue failed\tkey:%s\terr:%s", rq.key, err.Error()))
		return nil, false
	}
	return []byte(data), true
}

// Size 实现iface接口的队列长度方法
func (rq *RedisQueue) Size() int64 {
	len, err := rq.client.LLen(rq.key).Result()
	if err != nil {
		util.WriteLog(fmt.Sprintf("get queue size failed\tkey:%s\terr:%s", rq.key, err.Error()))
		return 0
	}
	return len
}
