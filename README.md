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

## 项目手册