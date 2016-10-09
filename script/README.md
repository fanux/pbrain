## 部署文档

### 系统架构图 依赖关系

```
rethinkdb      
    +------------->dface
                    ^
                    |
swarm---------------+
                    |
gnatsd              V
   +------------->pbrain------->plugins
                    ^
                    |
pgsql---------------+
```
### 镜像列表
* pgsql: postgres:latest
* gnatsd: gnatsd:latest
* rethinkdb: rethinkdb:latest
* pbrain: pbrain:latest
* dface: dface:latest

### 启动
使用docker-compose启动系统, 由于docker-compose的依赖关系不会等待进程启动成功，所以使用dependson会有问题，这里分开启动。

确保上一步执行完了再执行下一步, 在后台运行加-d参数
```
$ docker-compose -f docker-compose-dependson.yml up 
$ docker-compose -f docker-compose.yml up
$ docker-compose -f docker-compose-plugins up
```

### 需要优化的地方
由于dface的插件需要访问pbrain的地址。所以在 /go/static/app/shipyard.js 里配置 $rootScope.url

这个值在实际生产环境下需要改一下。

最好改成在dface的服务端获取这个值，服务端在启动时可以命令行指定。
