# Business(Worker) 被动处理


## 运行流程

### 1. `AsyncWebsocket`连接 `register注册发现`

### 2. 发送身份认证请求
### 3. （发生多次）等待 `register注册发现` 返回 `[]Gateway` 列表

### 4.  遍历过滤已连接的，剩下的通过 `AsyncWebsocket` 连接 `[]Gateway` 

4.1 如果  返回 `[]Gateway` 列表 包含有未知的即为`新扩容Gateway`

4.2 如果 返回 `[]Gateway` 列表 未包含 `已连接的Gateway` 即为被踢出集群 


### 若 `register注册发现`  断开，马上进行重连。但是不应影响现有的`[]Gateway`的连接保持

### 5. 连接`Gateway` 发送身份认证

### 6. 等待`Gateway` 被动消息：`要求认证`|| `用户连接`||`用户消息` || `用户断开`

### 7. 等待上方消息的处理，您的处理！