## Docker 镜像仓库搭建

### [Harbor](https://github.com/vmware/harbor)
VMware开源的docker镜像仓库, 可支持多镜像仓库之间的镜像同步

### 仓库分布
```
 192.168.86.170           192.168.86.106
 +----------+   sync      +----------+
 | registry | <---------> | registry |
 +----------+             +----------+
    ⬆
  develop
```
开发者提交docker镜像到仓库，仓库中建立两个项目

develop和deploy

开发者提交镜像到develop项目中

deploy项目同步到生产环境的仓库中。


### 仓库部署
安装docker-compose
```
$ pip install --upgrade pip && pip install docker-compose
```
```
$ git clone https://github.com/vmware/harbor
$ cd harbor
$ cd Deploy
```
编辑配置文件
```
$ vim harbor.cfg
```
修改 hostname参数为ip  或者主机域名(确保能连上)

```
$ ./prepare
Generated configuration file: ./config/ui/env
Generated configuration file: ./config/ui/app.conf
Generated configuration file: ./config/registry/config.yml
Generated configuration file: ./config/db/env

$ docker-compose up -d
```
浏览器访问对应域名或者ip即可

初始用户名/密码：admin/Harbor12345   (harbor.cfg文件中可配置)

离线部署请参考：[Harbor 离线部署](https://github.com/vmware/harbor/blob/master/docs/installation_guide.md)

重启：
```
$ docker-compose restart
```
删除镜像仓库容器
```
$ docker-compose rm
```

### 同步复制
假设上述170的服务器需要同步到106服务器上

* 浏览器进入170的仓库dashboard
* 登录后点击 管理员选项->系统管理->新建目标
* 新建目标窗口中填一个名称，目标url: http://192.168.86.106 用户名/密码: admin/Harbor12345 

配置项目策略

* 点击项目，新建一个项目如test
* 点击test进入项目 点击 复制->新增策略 
* 在目标设置中选定配置好的目标，点击确认
* 完成， 然后只要push到test项目的镜像就会自动同步到106的镜像仓库中。
