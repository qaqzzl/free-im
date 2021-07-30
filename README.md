### 简要介绍
```
即时通讯服务端
    系统分为3层，各服务通过grpc调用
        http服务      （登陆注册，好友列表，群列表，创建群，添加好友等。。。）
        im业务层       (消息处理）
        socket连接层   (长链接管理，消息接收，消息投递)
```

#### android run demo image
<img src="http://free-im-qn.qaqzz.com/docs/app1-1.jpg" width="25%"/><img src="http://free-im-qn.qaqzz.com/docs/app2-1.jpg" width="25%"/><img src="http://free-im-qn.qaqzz.com/docs/app3-1.jpg" width="25%"/><img src="http://free-im-qn.qaqzz.com/docs/app4-1.jpg" width="25%"/>

[Android 体验下载 ](https://www.pgyer.com/freeim)

[Android 项目地址](https://github.com/qaqzzl/free-im-android)


### 项目部署
1.拉取代码
```
git clone https://github.com/qaqzzl/free-im.git
# 创建MySQL数据库，导入sql文件
free-im/scripts/mysql.sql
```

2.配置文件
```
# 复制配置文件模板
cp ~/free-im/free.yaml.example ~/free-im/free.yaml
# 修改配置文件
vim ~/free-im/free.yaml
```

3.linux编译运行
```shell
cd ~/free-im

chmod a+x run.sh

./run.sh
```

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

### API接口文档
https://www.apizza.net/project/88dcdb7080c14030f5005d67132f5617/browse

