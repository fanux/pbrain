## Decider插件介绍 (业务驱动调度)
Decider （决策者） 之所以将此插件命名为这个是因为这个插件决定着容器的生杀大权。

一旦多个app都想伸缩，特别是“伸” 时，就需要权衡先满足谁。

甚至有时满足不了，需要杀死其它的app来满足某个优先级高的app。

决策者插件的策略配置文档就如同宪法一样，决策者的一切决策都是根据这个来的。

## 需要思考的问题
有很多问题需要思考，其实设计的重中之重就是策略文档的设计，这个设计出来了什么都好办了。

### App 优先级问题
起初准备给每个APP设置固定的优先级，这样会导致一个问题，就是可能把某一个app杀光了。
但是这往往是不想的结果。所以app的优先级应当随着数量的减少而增加，保证可控。

### Metrical（度量）问题
以什么纬度去决定调度，CPU？ 内存？ IO？ 其实都不合适，因为业务不同，伸缩的条件都不太一样，

所以决策者应该抽象出一个Metrcal(1~100)值，这个值怎么来的，插件也不知道，插件只知道对应的值来了怎么处理。

这个值估计只有业务知道，业务计算出这个值通知插件。

这个值可能是CPU，GPU，内存的使用情况，也可能是任务队列里的任务数量等。

## 策略文档设计
TODO Metrical最好是一个范围
```yaml
[
    {
        "App":"ats",
        "Specs":[
            {
                "Metrical":90,
                "Number":30,
                "Priority":50,
            },
            {
                "Metrical":60,
                "Number":20,
                "Priority":70,
            },
            {
                "Metrical":30,
                "Number":10,
                "Priority":90,
            },
            {
                "Metrical":10,
                "Number":5,
                "Priority":100,
            }
        ]
    },
    {
        "App":"hadoop",
        "Specs":[
            {
                "Metrical":90,
                "Number":20,
                "Priority":20,
            },
            {
                "Metrical":60,
                "Number":20,
                "Priority":40,
            },
            {
                "Metrical":30,
                "Number":20,
                "Priority":60,
            },
            {
                "Metrical":10,
                "Number":1,
                "Priority":100,
            }
        ]
    }
]
```
* `Metrical` 业务模块发来的一个度量值
* `Number`   收集到一个业务的消息，消息中的度量值对应的应该启动的实例数量 (TODO应该是伸缩的数量还是对应容器总数量？)
* `Priority` 优先级，很重要的一个参数，决定着当多个扩展实例的请求来了先满足哪个APP, 也决定着扩展请求来了资源不够用时释放哪些APP

`Number` 字段的问题:
假设业务任务队列中的任务数量（或者loadbalance的负载情况）与Metrical直接相关，任务多了处理不过来了业务就发送一个大的Metrical值，任务少了就发送一个小的Metrical值。所以照这样看Number应该是需要伸缩的数量，如`"Number":10` 代表新增加10个实例  `"Number":-5`代表减少5个实例。   这样带来两个问题：1 实例有可能被缩减完。 2  优先级怎么控制？ 给业务设置固定优先会导致每次只缩优先级低的，哪怕它很忙。   所以我觉得插件还是需要能够主动去获取业务负载情况的数据。


## 业务度量消息格式 
```json
{
    "App":"ats",
    "Metrical":80
}
```


---
华丽分割线

## 新的策略文档设计  （静态优先级）
```json
[
    {
        "App":"ats",
        "Priority":1,
        "MinNum":3,
        "Spec":[
            {
                "Metrical":[0, 20],
                "Number":-10,
            },
            {
                "Metrical":[20, 40],
                "Number":-5,
            },
            {
                "Metrical":[40, 60],
                "Number":0,
            },
            {
                "Metrical":[60, 80],
                "Number":5,
            },
            {
                "Metrical":[80, 100],
                "Number":10,
            }
        ]
    },
    {
        "App":"hadoop",
        "Priority":2,
        "MinNum":2,
        "Spec":[
            {
                "Metrical":[0, 20],
                "Number":-10,
            },
            {
                "Metrical":[20, 40],
                "Number":-5,
            },
            {
                "Metrical":[40, 60],
                "Number":0,
            },
            {
                "Metrical":[60, 80],
                "Number":5,
            },
            {
                "Metrical":[80, 100],
                "Number":10,
            }
        ]
    },
    {
        "App":"redis",
        "Priority":3,
        "MinNum":1,
        "Spec":[
            {
                "Metrical":[0, 20],
                "Number":-10,
            },
            {
                "Metrical":[20, 40],
                "Number":-5,
            },
            {
                "Metrical":[40, 60],
                "Number":0,
            },
            {
                "Metrical":[60, 80],
                "Number":5,
            },
            {
                "Metrical":[80, 100],
                "Number":10,
            }
        ]
    },
]
```
* `Priority`: [1-10] 1表示优先级高， 10表示优先级低, 高优先级的APP需要“伸”时，可释放低优先级的APP,  高优先级的APP处于空闲状态可让出一部分资源给低优先级的APP。

* `MinNum`: 实例最小的运行数量，到底线时不再进行“缩”操作。

业务仅仅上报一个值来描述负载情况：
```json
{
    "App":"ats",
    "Metrical":80    //数字越小表示资源越过剩，50表示均衡态，资源刚好够用，再大根据配置文件就有可能有扩展资源的动作了。
}
```

## 决策算法

一个“伸” 请求到来时

第一轮：挑出空闲资源APP，根据优先级进行释放, 先释放优先级低的，够用了即结束, 否则进入第二轮

第二轮：遍历优先级低于自己的APP，通过优先级和负载度量值给APP打分,得出各个APP让出多少资源出来用以满足高优先级的APP扩张的需求

> 低优先级的APP让高优先级APP释放资源，除非高优先级APP处于“空闲”状态，供过于求。

**打分**

Input(`ats`业务需要扩张10个实例)

Output(`hadoop`让出3个，`redis`让出7个)

算法体现出以下几点：
* 优先级越高让出资源越少
* 业务负载越高让出资源越少
* 让出资源后不得小于实例最小值数量（`MinNum`）
* (TODO要不要增加一个维度，当前容器数量)

1. 先给对应APP打分： score = 100 - Metrical + Priority*10
2. 需要扩张的总数量设为n， 某个app需要让出的数量为：n*score/(score1 + score2 + score3)
3. 计算出来的数值要与MinNum比较，如小于这个值就尽最大可能提供，也就是释放当前数量-MinNum个
4. 尽最大可能满足，有可能需要10个但是返回8个, 这个时候再进行一轮打分分配

---
华丽分割线

## 两个业务实现
前期的需求主要是两个业务之间的弹性调度。所以先针对两个业务的场景进行调度。
```json
[
    {
        "App":"ats",
        "Priority":1,
        "MinNum":3,
        "Spec":[
            {
                "Metrical":[0, 20],
                "Number":-10,
            },
            {
                "Metrical":[20, 40],
                "Number":-5,
            },
            {
                "Metrical":[40, 60],
                "Number":0,
            },
            {
                "Metrical":[60, 80],
                "Number":5,
            },
            {
                "Metrical":[80, 100],
                "Number":10,
            }
        ]
    },
    {
        "App":"hadoop",
        "Priority":2,
        "MinNum":2,
        "Spec":[
            {
                "Metrical":[0, 20],
                "Number":-10,
            },
            {
                "Metrical":[20, 40],
                "Number":-5,
            },
            {
                "Metrical":[40, 60],
                "Number":0,
            },
            {
                "Metrical":[60, 80],
                "Number":5,
            },
            {
                "Metrical":[80, 100],
                "Number":10,
            }
        ]
    },
]
```
现假设集群有总共有10个节点，初始状态，hadoop运行了5个，ats运行了5个，初始的Metrical值都为50, 下面来看看APP发来一些度量信息时会发生什么样的事。

> 收到消息

```json
{
    "App":"ats",
    "Metrical":55
}
```
查看配置表，55对应的Number值是0，也就是不需要扩张，所以什么都不做

> 收到消息

```json
{
    "App":"Hadoop",
    "Metrical":70
}
```
查看配置表，需要扩容5个实例。 但是“ats”的优先级高，而且度量值为55意味着没有空闲资源，所以不会给
Hadoop让出资源，且没有优先级更低的APP了，也什么都不做，仅报告一个扩张失败，解释失败原因。

> 收到消息

```json
{
    "App":"ats",
    "Metrical":75
}
```
这时表示“ats”需要扩容5个实例了，因为其优先级高，所以收缩hadoop业务来满足“ats”的需求，
不过hadoop仅有5个实例，而且配置了最小运行数量2
所以hadoop最多让出3个实例。此时让出3个实例给ats运行。

> 收到消息, ats高峰期过后，收到ats轻负载度量的消息

```json
{
    "App":"ats",
    "Metrical":30
}
```
表示ats很闲可以让出5个实例了，但是这个时候不马上让，等别人需要的时候再让。  

果然，因为hadoop被缩减了，一直处于高负载状态，发如下消息

```json
{
    "App":"Hadoop",
    "Metrical":90
}
```
这个时候发现高优先级的app有剩余资源，此时再去缩减ats。 

同样，如果缩减完了之后ats还是很空闲，hadoop还是很繁忙，那么同样根据收集到的消息进行再一轮的
伸缩，直到 ats的Metrical值在40～60之间。

最后，随着ats高峰期到来，进入一开始的循环。

PS: 上面配置文件的粒度很粗，实际情况可以配置更细粒度的度量值范围。


其中有很多细节问题需要考虑，如ats低负载时让出资源后要把Metrical值调到50（或者其他值）防止连续收到hadoop发来的高负载度量值，这样ats其实已经缩过了，只有等ats再次发来度量值时才能决定是否能

继续收缩。

## 业务发布消息格式
```go
COMMAND_APP_METRICAL = "app-metrical"

type Command struct {
	Command string
	Channel string 
	Body    string
}
```
json 格式为：
```json
{
    "Command":"app-metrical",      # 命令码，目前只有"app-metrical"
    "Channel":"plugin_decider",    # Channel是插件的名称，也是业务发布消息的关键(主题或通道)
    "Body":`{"App":"Ats","Metrical":90}`    # 注意Body是个json字符串
}
```
决策者插件会订阅“plugin_decider”通道，业务向这个通道publish上述格式消息即可.

## 另一种更简化的设计 (我们用一个新的插件实现这个，命名为plugin_justice)
这个比较针对ats和hadoop的业务场景，根据ats的负载情况决定集群中二者的运行数量
```json
{
    "App":"ats",
    "Sepc":[
        {
            "Metrical":[0-20],
            "Apps":[
                {
                    "App":"ats",
                    "Number":5,
                },
                {
                    "App":"hadoop:latest",
                    "Number":15,
                },
            ]
        },
        {
            "Metrical":[20-40],
            "Apps":[
                {
                    "App":"ats",
                    "Number":10,
                },
                {
                    "App":"hadoop:latest",
                    "Number":10,
                },
            ]
        },
        {
            "Metrical":[40-60],
            "Apps":[
                {
                    "App":"ats",
                    "Number":15,
                },
                {
                    "App":"hadoop:latest",
                    "Number":5,
                },
            ]
        }
    ]
}
```
可以看到配置中根据ats实例的数量决定各运行多少个ats，多少个hadoop
