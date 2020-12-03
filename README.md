#### android APP下载:  https://cdn.qaqzz.com/app-release.apk
#### android 源码:  https://github.com/qaqzzl/free-im-android

### 简要介绍
```
free-im是一个即时通讯服务器，代码全部使用golang完成。主要功能
    支持tcp，websocket接入(websocket暂未开发)
    离线消息同步
    单用户多设备同时在线
    单聊，群聊，以及超大群聊天场景(群聊未实现)
    登录注册
    修改用户资料
    好友增删改查,
    发布动态(支持图片, 视频)
```

### 目录结构
项目结构遵循 https://github.com/golang-standards/project-layout
```
cmd:          服务启动入口
config:       服务配置
internal:     每个服务私有代码
pkg:          服务共有代码
sql:          项目sql文件
test:         长连接测试脚本
```

### 使用技术：
```cgo
数据库：MySQL+Redis
应用数据格式：Protocol Buffers(暂时使用json)
通讯框架：Grpc  (暂未使用)
通讯协议: version(4) action(1) sequence-id(4) body-length(4) body-data
```

### 参考资料
    https://www.jianshu.com/p/9b58bb553cc0
    http://www.52im.net/thread-464-1-1.html
    // 58到家实时消息系统的协议设计等技术实践分享
    http://www.52im.net/thread-298-1-1.html
    https://blog.golang.org/migrating-to-go-modules
#####
    IM消息ID技术专题(一)：微信的海量IM聊天消息序列号生成实践算法原理篇
    http://www.52im.net/thread-1998-1-1.html
#####
    IM消息ID技术专题(三)：解密融云IM产品的聊天消息ID生成策略
    http://www.52im.net/thread-2747-1-1.html
#####
    参考项目
    https://github.com/alberliu/gim
#### 其他
```cgo
https://goproxy.io/
http://mirrors.aliyun.com/goproxy/
go list -m -json all
```


#### 命令笔记
```
grpc
go get -u github.com/golang/protobuf
go get -u github.com/golang/protobuf/protoc-gen-go
https://github.com/protocolbuffers/protobuf/releases
protoc --go_out=plugins=grpc:./pbs *.proto

```

#### 常见错误
```
exec: "gcc": executable file not found in %PATH%
下载gcc环境
https://jmeubank.github.io/tdm-gcc/download/

```

###思路

消息ID采用
![Image text](docs/message_id.png)

1. 消息储存设计
    redis 有序集合储存
        key: 会话ID
            score: 消息ID
            member: 消息

2. 离线消息设计
    规则: 
        设备上线自动推送
        每个群离线消息集合只储存最近1000条,
        每个单聊离校消息集合只储存最近5千条
        超过以上规则的可以用消息同步拉取
    存储:
        redis 有序集合储存
            key: 用户iD + 会话ID
                score: 消息ID
                member: 消息

3. 消息同步
    规则:
        只能根据通过会话ID跟消息ID, 同步当前会话ID 大于或小于 消息ID的记录
        
4. 消息回执

5. 客户端根据消息ID去重

(通过 1,2,3,4,5 保证消息不丢不重)