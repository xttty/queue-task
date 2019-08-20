package conf

import (
	"fmt"
	"queue-task/v1/util"
	"sync"

	"github.com/BurntSushi/toml"
)

// TomlConfig 总配置
type TomlConfig struct {
	Core  CoreConf
	Queue QueueConf
	Job   JobConf
}

// CoreConf 核心配置
type CoreConf struct {
	ConfPath string
}

// QueueConf 队列配置
type QueueConf struct {
	Redis map[string]RedisQueueConf
}

// RedisQueueConf redis队列连接配置
type RedisQueueConf struct {
	Addr     string
	Password string
	DB       int
}

// JobConf 任务配置
type JobConf struct {
	Default DefaultJobConf
}

// DefaultJobConf 默认任务配置
type DefaultJobConf struct {
	WorkersCnt int
}

// Config 总配置
var Config *TomlConfig
var once sync.Once

// Init 初始化
func Init(confPath string) {
	once.Do(func() {
		Config = &TomlConfig{
			Core: CoreConf{
				ConfPath: confPath,
			},
		}
		qConfPath := Config.Core.ConfPath + "/queueconf.toml"
		_, err := toml.DecodeFile(qConfPath, &(Config.Queue))
		if err != nil {
			util.WriteLog(fmt.Sprintf("queue config load failed, err: %s", err.Error()))
		}
		jConfPath := Config.Core.ConfPath + "/jobconf.toml"
		_, err = toml.DecodeFile(jConfPath, &(Config.Job))
		if err != nil {
			util.WriteLog(fmt.Sprintf("job config load failed, err: %s", err.Error()))
		}
	})
}
