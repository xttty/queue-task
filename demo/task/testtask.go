package task

import (
	"encoding/json"
	"fmt"
	"queue-task/demo/conf"
	"queue-task/v1/iface"
	"queue-task/v1/job"
	"queue-task/v1/msg"
	"queue-task/v1/queue"
)

// CreateTestJob 新建测试队列任务
func CreateTestJob() CreateJobFunc {
	return func() iface.IJob {
		// 初始化一个队列
		q := queue.NewKaproxy(conf.Config.Queue.Kaproxy["default"], TestKaproxyTopic, TestKaproxyGroup)
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
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(message)
}
