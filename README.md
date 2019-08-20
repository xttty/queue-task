## queue\-task
### 简介
**queue\-task**是一个轻量级方便扩展的队列调度框架。
### 使用
1. 在项目目录使用`dep ensure`加载配置
2. 详细使用参考demo文件夹中的使用案例。
demo包含以下内容：
* daemon 守护进程（因为不同业务，守护进程也不同，所以不放在框架内部）
  * 任务注册（将所有需要启动的任务都放在方法中）
  * 任务启动，在任务注册后启动任务
* conf 配置文件，如果需要扩展配置需要修改conf/config.go文件
* task 包含具体的业务代码以及业务注册

### 核心功能
框架核心文件在v1目录下，主要包含：
* iface 框架基本接口定义，包括（**任务接口**、**队列接口**、**消息接口**）
* job 具体实现任务接口的任务模板，框架本身自带两种任务模板：
  * basejob 基础的任务模板，仅具体实现了业务方法注册
  * defaultjob 默认任务模板，实现了任务的并发
* queue 任务队列的抽象，隐藏了具体队列
* msg 消息的封装
* util 框架公用方法封装
* conf 核心配置（队列配置、job配置、核心配置）