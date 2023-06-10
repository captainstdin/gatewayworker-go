# gatewayWorker-Go

## 简介：
这是一个从gatewayworker模型中得到的启发，使用`golang-1.20`实现。

## 此项目 包含组件
+ `business` 对照gatewayworker(workerman)
+ `register` 对照gatewayworker(workerman)
+ `gateway` 对照 gatewayworker(workerman)
+ `sdk` 对照 GatewayClient(workerman)

补充：
1. `GatewayClient` 就是原来的其他PHP项目主动推送
https://www.workerman.net/doc/gateway-worker/push-in-other-project.html

2. sdk组件包含一个`SdkClient`，提供主动推送能力；一个`WebServer`，提供`WebAdmin`管理界面

## 对比gatewayWorker(workerman)改进
+ 全面支持IPv6和IPv4(包括用户与`gateway`的链路，也包括`gateway`,`reigster`,`business`之间的通讯)

+ 组件之间通讯全面使用ws+tls(可选)+签名校验（原workerman用的是`secertKey`校验权限）

+ 支持Sdk更多，增加了Http服务器 `复用` Gateway的端口，此服务器内置web监控

+ 统一`服务端口`和`组件内部通讯端口` 仅需`1个`端口对外，即(http:// ws:// 使用相同端口)。

补充：
+ 原`Gateway`(workerman)使用 一个服务端口，另外N个进程`2900`,`2901`....也需要对外使用，主要是给`GatewaySdk`和`Business`连接通讯使用

+ 由于使用`协程非阻塞`，所以在推送的时候可以`并发推送`，原`GatewayClient`(workerman)是阻塞式for遍历推送

## 项目手册


## 部署模式
1. (难度：★☆☆)服务器单节点部署


### 1.1 构建

(可选)github Action构建


本地机器构建

```
#如果需要大陆镜像加速
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GO111MODULE=on
go mod tidy 
go build -o myapp .
```


2. (难度：★★☆) 服务器集群部署

> 详见手册， 支持(Docker部署)


3. (难度：★★★) Serverless部署`gateway`与`business`，一台非高可用的`公网服务器`用于`组件内部广播`

> 详见手册， 支持(Docker部署)
