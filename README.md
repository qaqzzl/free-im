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
go list -m -json all
```
