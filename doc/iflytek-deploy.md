## UI 组件：
* Dface http://172.16.162.4:8888/ (admin/shipyard) 可在上面看到集群节点情况，容器列表，镜像列表，以及对容器的一些常用操作
* Harbor http://172.16.162.10/ (admin/Harbor12345) 镜像仓库，develop项目低下可看到仓库中的镜像列表。
* Gogs http://172.16.162.4:3000 git仓库，可用可不用。。。

## 镜像构建
将包含Dockerfile的目录拷贝到172.16.162.4服务器上的任意目录执行build命令，以dface为例：
```
docker build -t dface:latest .
```
## 提交镜像仓库
```
docker tag dface:latest reg.iflytek.com/develop/dface:latest
docker push reg.iflytek.com/develop/dface:latest
```
注意事项：

* `develop`是harbor中建立项目名称，harbor界面上可看到。
* 若文件有更新重新build镜像并push，镜像名称一样会覆盖之前仓库中的镜像，如不想覆盖请使用新的tag。

## 执行compose
```
version: '2'
services:
    iat_162_25:
        # 指定容器名称
        container_name: iat_162_25
        # 指定容器使用的镜像,更新新版本时需要改动
        image: reg.iflytek.com/develop/iat:latest
        # 节点过滤器，可以使该容器运行在有hostname=yjybj-162-025的节点上
        environment:
            - "constraint:hostname==yjybj-162-025"
        # 容器启动命令
        command: sh run_iatserver.sh 172.16.162.25:9090
        # 容器需要挂载的目录
        volumes:
           - /disk0/iat/resource:/root/resource
           - /disk0/iat/data:/root/data
           - /disk0/iat/conf/iat_server.ini:/root/sivs_run/bin/iat_server.ini
        # 容器的网络模式，统一使用host模式
        network_mode: "host"
```

启动容器:
```
docker-compose -H tcp://172.16.162.4:4000 up -d
```
* -H 指定swarm manager的地址，统一为`tcp://172.16.162.4:4000`, 运行容器的请求会发送给swarm manager，swarm manager转发给指定主机的docker engine
* -d 在后台启动容器
* -f 可指定compose.yml文件

## 删除容器
有几种方式：
```
docker-compose -H tcp://172.16.162.4:4000 down
```

在Dface界面上点击Destroy删除容器

或者到宿主节点使用docker命令删除

## 清除目标机缓存镜像

在更新业务时如果镜像名称没变的话，宿主机上缓存的镜像不会因为仓库的镜像更新了而自动更新。  需要删除宿主机(运行容器的目标机)上的缓存.

可以在Dface界面的IMAGES菜单下找到对应镜像删除。

或者到目标机上执行 `docker rmi $(镜像名)`命令

## 查看容器日志
在Dface上点击容器详情，可看到容器的logs信息，这些logs是业务打印到标准输出标准错误的日志。

如果业务挂载出日志文件需要到对应主机目录下查看日志文件

## 进入容器
如有特殊需求想进入容器内部需要到容器运行的机器上执行:
```
docker exec -it $(容器名) /bin/bash
```
不想到宿主机上去也行：
```
docker -H tcp://172.16.162.4:4000 exec -it iat_162_25 /bin/bash
```
-H 指定swarm manager的地址
