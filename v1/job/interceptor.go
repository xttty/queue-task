package job

import "queue-task/v1/iface"

// WorkInterceptor  job运行中间件
// 拦截器一般用于数据统计、日志等
// 拦截器需要显示调用 handle(value) 来继续业务处理流程
// handle可以在拦截器代码的任何位置以此来控制拦截器逻辑相对于业务处理的位置
type WorkInterceptor func(value []byte, handle iface.JobHandle)

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
		for i := 0; i < n; i++ {
			chainHandle = chainHandleFunc(workInts[i], chainHandle)
		}
		chainHandle(value)
	}
}
