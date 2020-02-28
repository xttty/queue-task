## queue\-task
### 简介
**queue\-task**是一个轻量级、易扩展、支持分布式的队列调度框架。
### Required
- go version >= 1.12
- 开启go mod
### quick start
具体使用案例可看`./examples`文件
#### 单机任务创建
```go
// 创建队列
// 队列类型可以从框架中选择（目前支持redis\kafka）也可以自定义
q := queue.NewRedisQueue(&queue.RedisQueueOptions{
  Addr: "redis.dev:6379",
  Key:  "debug",
})
// 使用任务模板
testJob := job.NewDefaultJob("test", q, 10)
// 注册业务handle方法
testJob.RegisterHandleFunc(func(data []byte) {
  // ...
})
```
#### 分布式任务创建
```go
// 创建队列
q := queue.NewRedisQueue(&queue.RedisQueueOptions{
  Addr: "redis.dev:6379",
  Key:  "debug",
})
// 分布式控制策略
strategy := func() int {
  var jobCnt int
  // ...
  return jobCnt
}
// 创建分布式队列
job := job.NewDistributedJob("debug", q, job.SetCircleTime(5*time.Second), job.SetWorkerStrategy(strategy))
```
### 定制化
#### 定制消息
实现`./v1/iface/imsg.go`中`IMessage`接口，即可自定义消息。
#### 定制队列
实现`./v1/iface/iqueue.go`中定义的`IQueue`接口即可自定义消息队列，使用方式和原生队列无差别。
#### 定制任务
实现`./v1/iface/ijob.go`中定义的`IJob`接口，即可自定义任务，或者“继承”框架内部的`BaseJob`也可以自定义任务。

### 框架核心
框架核心文件在v1目录下，主要包含：
* iface 框架基本接口定义，包括（**任务接口**、**队列接口**、**消息接口**）
* job 具体实现任务接口的任务模板，框架本身自带两种任务模板：
  * basejob 基础的任务模板，仅具体实现了业务方法注册
  * defaultjob 默认任务模板，实现了任务的并发
  * distributedjob 分布式任务模板，实现了在分布式环境下的任务调控
* queue 任务队列的抽象，隐藏了具体队列，目前支持的队列：
  * redis队列
  * kafka消息队列
* msg 消息的封装
* task 所有任务通过此包进行统一管理
* util 框架公用方法封装

### Demo环境搭建
进入项目目录执行`docker-compose up -d`

详情请看**docker-compose.yml**文件