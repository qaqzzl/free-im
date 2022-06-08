### 简要介绍
```
即时通讯服务端
    系统分为3层，各服务通过grpc调用
        http服务
        im业务层
        socket连接层
```

#### 安卓客户端效果图
<img src="http://free-im-qn.qaqzz.com/docs/app1-1.jpg" width="25%"/><img src="http://free-im-qn.qaqzz.com/docs/app2-1.jpg" width="25%"/><img src="http://free-im-qn.qaqzz.com/docs/app3-1.jpg" width="25%"/><img src="http://free-im-qn.qaqzz.com/docs/app4-1.jpg" width="25%"/>

安卓APP下载：[https://www.pgyer.com/freeim](https://www.pgyer.com/freeim)

GitHub地址：[https://github.com/qaqzzl/free-im-android](https://github.com/qaqzzl/free-im-android)


### 项目部署 (Linux)

```
安装软件：
mysql5.7, redis6.2

# 拉取代码
git clone https://github.com/qaqzzl/free-im.git
# 创建MySQL数据库，导入sql文件
free-im/scripts/mysql.sql

# 复制配置文件模板
cp ~/free-im/free.yaml.example ~/free-im/free.yaml
# 修改配置文件
vim ~/free-im/free.yaml
# 下载项目依赖
go mod download

# linux编译运行
cd ~/free-im
chmod a+x run.sh
./run.sh
```


### 核心流程时序图
> 登陆

<img src="http://free-im-qn.qaqzz.com/docs/login.png" width="70%"/>

[comment]: <> (> 消息同步)

[comment]: <> (<img src="http://free-im-qn.qaqzz.com/docs/message_sync.png" width="70%"/>)

> 发送消息(单聊)

<img src="http://free-im-qn.qaqzz.com/docs/message_send.png" width="70%"/>

[comment]: <> (> 发送消息&#40;群聊&#41;)

[comment]: <> (<img src="http://free-im-qn.qaqzz.com/docs/message_group_send.png" width="70%"/>)



### 常见错误
golang.org 包拉不下来
```
export GOPROXY=https://mirrors.aliyun.com/goproxy/
```

Windows运行项目需要gcc环境
```
exec: "gcc": executable file not found in %PATH%
下载gcc环境
https://jmeubank.github.io/tdm-gcc/download/
```



