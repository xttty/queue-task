package msg

import "encoding/json"

// BaseMsg 基础消息
type BaseMsg struct {
	Data H
}

// H 类似php的数组使用姿势
// 	msg.H{
// 		"abc": 123,
// 		"xty": msg.H{
// 			"x": 1,
// 			"y": 2,
// 			"map": msg.H{
// 				"123": "234",
// 			},
// 		},
// 	}
type H map[string]interface{}

// ToJSON 实现的消息接口
func (msg *BaseMsg) ToJSON() ([]byte, error) {
	v, err := json.Marshal(msg)
	return v, err
}
