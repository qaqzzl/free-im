### 简要介绍
free-im是一个即时通讯服务器，代码全部使用golang完成。主要功能
1.支持tcp，websocket接入(websocket暂未开发)
2.离线消息同步
3.多业务接入(未实现)
4.单用户多设备同时在线
5.单聊，群聊，以及超大群聊天场景
### 使用技术：
数据库：Mysql+Redis
长连接通讯协议：Protocol Buffers



### 参考资料
https://www.jianshu.com/p/9b58bb553cc0
http://www.52im.net/thread-464-1-1.html

https://blog.golang.org/migrating-to-go-modules

https://goproxy.io/

go list -m -json all



ordinary  普通聊天
group     群聊