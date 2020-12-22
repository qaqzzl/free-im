## 客户端
android APK下载:  https://cdn.qaqzz.com/app-release.apk
    
android 项目地址:  https://github.com/qaqzzl/free-im-android
### 简要介绍
```
    即时通讯服务器，代码全部使用golang完成。主要功能
    支持tcp，websocket接入(websocket暂未开发)
    单用户多设备同时在线
    单聊，群聊，以及超大群聊天场景
    登录注册
    修改用户资料
    好友增删改查
    发布动态
```

### 目录结构
项目结构遵循 https://github.com/golang-standards/project-layout


### 使用技术：
```cgo
数据库：MySQL+Redis
应用数据格式：Protocol Buffers(暂时使用json,方便开发调试)
通讯框架：Grpc
通讯协议: 自定义
version(4) action(1) sequence-id(4) body-length(4) body-data
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



#### 笔记
```
mod厂库代理
https://goproxy.io/
http://mirrors.aliyun.com/goproxy/
go list -m -json all

grpc安装
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
```text
消息储存设计, 
    (写扩散) 群聊|单聊
    redis 有序集合储存
        key: 用户ID
            score: 消息ID
            member: 消息
    (读扩散) 聊天室:超大群聊
    redis 有序集合储存
        key: 会话ID
            score: 消息ID
            member: 消息


消息同步
    规则:
        通过消息ID, 同步消息记录


客户端根据消息ID去重


客服端负责消息重发, 保证服务端无状态性

********客户端A -----message1 to B----->   连接层   --------------> 逻辑层
********客户端B <----message1-----------   连接层   <-------------- 逻辑层
********客户端B -----message1 ack to A->   连接层                   逻辑层
********客户端A <----message1 ack-------   连接层                   逻辑层
```