package conf_test

import (
	"queue-task/demo/conf"
	"testing"
)

func TestInit(t *testing.T) {
	conf.Init()
	t.Log(conf.ProjectConfig)
}
