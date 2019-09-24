package task

import (
	"encoding/json"
	"fmt"
	"queue-task/v1/conf"
	"queue-task/v1/iface"
	"queue-task/v1/job"
	"queue-task/v1/msg"
	"queue-task/v1/queue"
	coretask "queue-task/v1/task"
	"queue-task/v1/util"
	"time"
)

// CreateTestJob 新建测试队列任务
func CreateTestJob() coretask.CreateJobFunc {
	return func() iface.IJob {
		// 初始化一个队列
		queueConf := conf.Config.Queue.Redis["test"]
		q := queue.NewRedisQueue(&queueConf, TestRedisJobKey)
		// 生成job
		j := job.NewDefaultJob("test", q, &conf.Config.Job.Default)
		// 注册业务回调方法
		j.RegisterHandleFunc(TestPerform)
		return j
	}
}

// TestPerform 测试业务代码
func TestPerform(data []byte) {
	var message msg.BaseMsg
	err := json.Unmarshal(data, &message)
	if err == nil {
		tempTime := message.Data["date"].(float64)
		time := int(tempTime)
		fmt.Println(time)
	} else {
		util.WriteLog(err.Error())
	}
}

func TestSendMsg() {
	createFunc := CreateTestJob()
	testJob := createFunc()
	for i := 0; i < 10; i++ {
		message := &msg.BaseMsg{
			Data: msg.H{
				"date": time.Now().Second(),
			},
		}
		testJob.Send(message)
		time.Sleep(time.Second)
	}
}
