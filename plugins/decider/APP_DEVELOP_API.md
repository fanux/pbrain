## 业务驱动概览

```
       app            plugin 
        |ats metrical 70 | 
        |--------------->|
        |     1          |
        |                |
        |  host ip list  |
        |<---------------|
     +--|     2          |
stop |  |                |
the  |  |                |
hosts|__|stop apps       |
apps    |--------------->|
        |emit scale down |
        | request        |
        |     3          |
```

可以看到，业务与decider插件主要有以上三个交互

* 1: 给插件发送负载度量值
* 2: 接收哪些节点上的APP将被释放
* 3: 处理完对应节点上的APP，通知插件可以进行释放操作了。

所以业务开发主要关心上面三个流程

## 消息队列
APP与插件通过消息队列相互通信。使用[nats](https://nats.io/)消息队列

[nats java client](https://github.com/nats-io/jnats)

## 消息格式
```go
type Command struct {
	Command string
	Channel string 
	Body    string
}
```

APP订阅自己（需要伸缩的APP镜像名称）

decider插件会订阅 `plugin-decider` 

#### 给插件发送负载度量值
业务publish一条消息
```json
{
    "Command":"app-metrical",      # 命令码
    "Channel":"plugin_decider",    # Channel是插件的名称，也是业务发布消息的关键(主题或通道)
    "Body":`{"App":"Ats","Metrical":90}`    # 注意Body是个json字符串
}
```
* App : 对应app镜像名称
* Metrical : 负载度量值，0 ～ 100 ，50表示负载平衡的状态，可以web界面上配置。

#### 接收哪些节点上的APP将被释放
插件publish一条消息，Channel是APP的镜像名称

命令码是`app-scale`

业务订阅如`hadoop:latest`关键字

```json
{
    "Command":"app-scale",      # 命令码
    "Channel":"hadoop:latest", 
    "Body":`{
        "App":"hadoop:latest",
        "Number":-2,            # 释放两个实例
        "MinNumber":1,          # APP无需关心
        "Hosts":[               # 释放哪几个节点上的实例
            "192.168.0.2",  
            "192.168.0.3",
        ],
        "ScaleUp":"ats:latest"  # 此字段APP不用关心，按原样发回给decider插件即可 
    }`  # 注意Body是个json字符串
}
```
注：这里面有些字段APP无需关心却也发送过来了，原因是decider插件是无状态的，这样可以保持一致性。如同一个原子操作不可分割

#### 通知插件可以进行释放操作了
业务publish一条消息, 和上个步骤消息格式一样，只是Channel的值变成了`plugin_decider`

```json
{
    "Command":"app-scale",      # 命令码
    "Channel":"plugin_decider", 
    "Body":`{
        "App":"hadoop:latest",
        "Number":-2,         
        "MinNumber":1,
        "Hosts":[
            "192.168.0.2",  
            "192.168.0.3",
        ],
        "ScaleUp":"ats:latest"  
    }`  # 注意Body是个json字符串
}
```
