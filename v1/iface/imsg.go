package iface

// IMessage 消息接口
type IMessage interface {
	ToJSON() ([]byte, error)
}
