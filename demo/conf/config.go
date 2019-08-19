// Package conf 配置文件包，所有配置文件使用toml语法编写
package conf

import (
	"queue-task/v1/util"
	"sync"

	"github.com/BurntSushi/toml"
)

// TomlConfig 项目配置
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
	Kaproxy map[string]KaproxyQueueConf
}

// KaproxyQueueConf kafka队列配置
type KaproxyQueueConf struct {
	Host  string
	Port  int
	Token string
}

// JobConf 任务配置
type JobConf struct {
	Default DefaultJobConf
}

// DefaultJobConf kafka任务配置
type DefaultJobConf struct {
	WorkersCnt int
}

// Config 配置实例
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
			util.WriteLog(err.Error())
		}
		jConfPath := Config.Core.ConfPath + "/jobconf.toml"
		_, err = toml.DecodeFile(jConfPath, &(Config.Job))
		if err != nil {
			util.WriteLog(err.Error())
		}
	})
}
