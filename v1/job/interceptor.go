package job

import "queue-task/v1/iface"

// WorkInterceptor  job运行拦截器
// 拦截器一般用于数据统计、日志等
// 拦截器需要显式调用 handle(value) 来继续业务处理流程
// handle可以在拦截器代码的任何位置以此来控制拦截器逻辑相对于业务处理的位置
type WorkInterceptor func(value []byte, handle iface.JobHandle)

// SendInterceptor  发送拦截器
// 发送拦截器用法和job运行拦截器一样
type SendInterceptor func(msg iface.IMessage, handle iface.SendHandle)

// ChainWorkInterceptor 链式work拦截器生成组件
func ChainWorkInterceptor(workInts ...WorkInterceptor) WorkInterceptor {
	n := len(workInts)

	return func(value []byte, handle iface.JobHandle) {
		chainHandleFunc := func(currentWorkInt WorkInterceptor, currentHandle iface.JobHandle) iface.JobHandle {
			return func(currentVal []byte) {
				currentWorkInt(currentVal, currentHandle)
			}
		}
		chainHandle := handle
		for i := n - 1; i >= 0; i-- {
			chainHandle = chainHandleFunc(workInts[i], chainHandle)
		}
		chainHandle(value)
	}
}

// ChainSendInterceptor 链式send拦截器生成组件
func ChainSendInterceptor(sendInts ...SendInterceptor) SendInterceptor {
	n := len(sendInts)

	return func(msg iface.IMessage, handle iface.SendHandle) {
		chainHandleFunc := func(currentInt SendInterceptor, currentHandle iface.SendHandle) iface.SendHandle {
			return func(currentMsg iface.IMessage) {
				currentInt(currentMsg, currentHandle)
			}
		}

		chainHandle := handle
		for i := n - 1; i >= 0; i-- {
			chainHandle = chainHandleFunc(sendInts[i], chainHandle)
		}
		chainHandle(msg)
	}
}
