## DFace 插件开发指南
使用pbrain框架可以帮助开发者快速开发DFace插件

## pbrain介绍
pbrain是一个go语言框架，里面包含了一些通用的与plugin manager交互使用的库，和一些定义好的结构体与关键字等。

pbrain抽象出了一些方法，具体的插件只需要实现对应接口，pbrain会调用这些接口，这使得开发业务插件不用关心多协程和

消息机制如何实现的，大大简化了插件开发的难度。

pbrain 一套代码中可包含多个插件，每个插件可通过命令行独立运行。

如：
```
$ pbrain pipeline
          插件名
```
表示运行pipeline插件，代码里的其它插件不会被执行 


```
▾ cmd/
    pipeline.go             // 一个具体的插件
    plugin.go               // 主逻辑，启动插件和启动协程监听plugin manager发来的事件
    root.go
▾ common/
    api.go                  // 插件通过http与manager通信的实现
    client.go               // 插件与manager通信的接口抽象，意味着可以不通过http通过，可扩展其他实现如protobuf
    command_test.go
    define.go               // 定义了通用的结构体，如插件，策略等还有通用的关键字定义
    interface.go            // 插件接口抽象，每个插件需要实现的接口，有些通用的接口插件基类实现
    mq.go                   // 消息队列的封装，用的是nats消息队列
    plugin_base.go          // 插件基类，所有插件聚合这个基类插件，就可不用重复实现一些通用功能，如获取插件名
    README.md
▾ plugins/pipeline/
    pipeline.go             // 具体的插件，实现interface里的接口
  LICENSE
  main.go
```

## 执行流程
pbrain是一个标准的cobra工程，当我们需要开发一个新的插件时，到pbrain目录执行(以pipeline插件为例)：
```bash
$ cobra add pipeline
```
这样就自动会在cmd目录下面生成pipeline.go文件

可在这个文件中处理命令行参数，然后运行`RunPlugin()`

```go
RunPlugin(&pipeline.Pipeline{common.GetBasePlugin(ManagerHost, ManagerPort, pipeline.PLUGIN_NAME), nil})  
```
```go
func RunPlugin(plugin common.PluginInterface) 
```
可以看到RunPlugin接受一个接口，我们传入的是一个具体的插件Pipeline，Pipeline聚合了基插件

RunPlugin 主要逻辑：

1. 创建一个channel用于主协程和监听事件的协程通信。

2. 主协程启动插件

3. 启动监听事件协程

4. 主协程循环等待事件发生，在事件发生时作对应处理

5. 监听协程中订阅了插件名通道，用于监听plugin manager发布的消息, 然后报告给主协程


RunPlugin里规定了，收到什么命令启动插件，停止插件，启动策略，停止策略，更新文档

事件发生时调用对应接口。

所以开发一个插件仅需要实现这些接口即可

## Pipeline插件介绍
```go
type Pipeline struct {
    *common.PluginBase
    // the key is strategy name
    Jobs map[string]Job
}
```
这里聚合和PluginBase，PluginBase包含了一些能用的插件信息

```go
type PluginBase struct {
    Client     Client
    Plugin     Plugin
    Strategies []Strategy
}
```
Client可与plugin manager通信，如获取插件信息，获取策略等

Plugin用于保存插件信息

Strategies用于保存策略信息

PluginBase还提供一些通用的方法，如初始化插件，获取插件名， 获取策略等

> Pipeline Start():
启动定时器，遍历所有插件的策略，将配置信息加入到任务中


