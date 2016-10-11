## Dockerfile 快速入门

### Why Dockerfile
docker镜像有几种构建方式，比如进入到容器中，安装好我们需要的东西，启动工程，然后基于这个容器导出镜像。

这样做不方便的地方是如果我们修改了代码或者配置什么的，就需要把重复的工作再来一遍，麻烦。

dockerfile指定了镜像构建的规则，这样就避免了手动构建的麻烦。

此外dockerfile的好处就是我们可以看到镜像构建的过程，镜像有什么问题可以找出来。

### 简单事例
pbrain是go语言的集群调度插件系统，我是在本地（开发环境）开发，在容器中运行的。

> 目录结构

```
▾ pbrain/
  ▸ cmd/
  ▸ common/
  ▸ doc/
  ▸ Godeps/
  ▸ manager/
  ▸ plugins/
  ▸ script/
  ▸ vendor/
    Dockerfile
    LICENSE
    main.go
    README.md
```
在工程主目录下创建一个Dockerfile文件

> Dockerfile

```bash
FROM golang:latest                            # 选择一个go的基础镜像，这样编译运行的环境就有了

COPY . /go/src/github.com/fanux/pbrain/       # 把代码拷贝到容器中

RUN go get github.com/tools/godep && \        # 在容器中编译安装
    cd /go/src/github.com/fanux/pbrain/ && \
    godep go install

CMD pbrain --help                             # 容器启动时执行的命令
```
Dockerfile非常的简单, 需要注意的地方有：

1. java 这种build once run any where的就可以不用拷源码到容器中了，直接拷贝可执行jar包，配置，执行脚本什么的，也就是运行需要的依赖。

2. 推荐把配置文件或者配置目录也打入镜像中，除非容器在不同的机器上运行时配置不同，如果是那样的话，可在宿主机上配置然后运行时作磁盘映射，尽量不要这样做。

3. RUN 命令可以在容器中执行一条linux命令。

4. CMD 命令可以在启动容器时被覆盖，如 我们启动容器
```bash
$ docker run pbrain:latest pbrain manager
```
这样`pbran manager` 这条命令就会覆盖镜像中的 `pbrain --help` 命令

> build镜像

在有Dockerfile的目录下执行：
```bash
$ docker build -t pbrain:latest .
$ docker images
```
（不要忘记后面的一点）
这样使用docker images命令就可以看到新的镜像 `pbrain:latest`

> 其它Dockerfile命令

其它命令如需要使用或者有不明白的地方可联系我

* MAINTAINER 用来指定镜像创建者信息
* ENTRYPOINT 设置container启动时执行的操作, 和CMD很像，但是不会被覆盖
* USER 设置容器启动时用户，默认是root用户。
* EXPOSE 暴露端口，我们共享主机网络，不用这个
* ENV 设置环境变量，这个可能需要用到, 当然也可以在容器运行时指定环境变量
* ADD 和COPY很像，就用COPY就可以了
* VOLUME 指定挂载点 挂载本机的目录，配置文件尽量不要挂载，数据输出可以，可在运行时指定
* WORKDIR 切换目录 切换工作目录，这个比较有用, 为了方便可以把运行时文件拷贝到此目录

### 提交镜像到镜像仓库
假设我们仓库地址是：`192.168.86.106:5000`

我们需要提交的镜像是`pbrain:latest`
```bash
$ docker tag pbrain:latest 192.168.86.106:5000/pbrain:latest
$ docker push 192.168.86.106:5000/pbrain:latest
```

### Docker compose

虽然有了Dockerfile 我们可以很方便的构建镜像，但是还有个问题，如何运行这些镜像？

docker compose帮我们解决了这个问题

创建一个docker-compose.yml文件，内容如下：
```yaml
version: '2'
services:
    pbrain:
      container_name: pbrain                                     # 指定容器运行名称
      command: pbrain manager -o http://192.168.86.106:8888      # 容器启动命令
      network_mode: "host"                                       # 指定网络模式
      image: 192.168.86.106:5000/pbrain:latest                   # 容器镜像
```
执行:
```
$ docker-compose up
```
其它参数：
* -f 指定docker compose文件，默认是docker-compose.yml
* -d daemon模式运行容器
* -H 指定docker engine 地址j

注：

* 一个compose文件可定义多个容器的运行，也可以定义容器间依赖关系。
* 不推荐使用compose定义依赖关系，因为在实践时发现compose不会等依赖容器进程真正运行起来才去运行被依赖容器。 所以有依赖的话分开，写多个compose文件，多up几次。 
* 开发人员需要提供镜像和compose配置文件。
