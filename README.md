### 简要介绍
```
    即时通讯服务端, 主要功能
    支持tcp，websocket接入
    单用户多设备同时在线
    单聊，群聊，以及超大群聊天场景
```

## 安卓客户端
[Android 安装包下载](https://cdn.qaqzz.com/app-free-release-v1.apk)

[Android 项目地址](https://github.com/qaqzzl/free-im-android)


### 部署
安装
```
# 拉取
git clone https://github.com/qaqzzl/free-im.git
# 创建MySQL数据库，导入sql文件
free-im/scripts/mysql.sql
```

配置文件
```
# 复制配置文件模板
cp ~/free-im/free.yaml.example ~/free-im/free.yaml
# 修改配置文件
vim ~/free-im/free.yaml
```

linux编译运行
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

### android demo image
<img src="http://free-im-qn.qaqzz.com/docs/app1.png" width="25%"/><img src="http://free-im-qn.qaqzz.com/docs/app2.png" width="25%"/><img src="http://free-im-qn.qaqzz.com/docs/app3.png" width="25%"/><img src="http://free-im-qn.qaqzz.com/docs/app4.png" width="25%"/>
