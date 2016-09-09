
## 更简化的业务驱动调度设计 plugin_justice
这个比较针对ats和hadoop的业务场景，根据ats的负载情况决定集群中二者的运行数量

这种设计比较适合有核心业务的场景，其他业务都围绕着核心业务转

这个设计也很有用，因为伸缩某个核心业务时可能对应的一些服务的业务也得进行伸缩。
通过这个配置就可以做到。

缺点是目前不支持多种业务驱动调度，多种业务就涉及到优先级的问题了，decider插件会去做这个事。
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
