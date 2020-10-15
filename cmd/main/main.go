package main

import (
	"flag"
	"free-im/configs"
	"free-im/internal/main/api"
	"log"
	"net/http"
	"os"
)



func main() {
	// 配置文件
	ConfPath := flag.String("cpath", "./config.conf", "config file")
	config.InitConfig(*ConfPath)

	init_http()
}


func init() {
	file := "./" +"message"+ ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}


func init_http() {
	http.HandleFunc("/login", api.PhoneLogin)		// 手机号登录 / 注册
	http.HandleFunc("/login/send.login.sms", api.SendLoginSms)		// 发送登录手机号验证码
	http.HandleFunc("/user/member.info", api.GetMemberInfo)		// 获取会员信息
	http.HandleFunc("/user/update.member.info", api.UpdateMemberInfo)		// 修改会员信息
	http.HandleFunc("/user/others.home.info", api.OthersHomeInfo)		// 获取用户基本信息(他人主页)
	http.HandleFunc("/search/friend", api.SearchMember)		// 搜索好友
	http.HandleFunc("/user/add.friend", api.AddFriend)		// 添加好友
	http.HandleFunc("/user/del.friend", api.DelFriend)		// 删除好友
	http.HandleFunc("/user/friend.apply.list", api.FriendApplyList)		// 好友申请列表
	http.HandleFunc("/user/friend.apply.action", api.FriendApplyAction)		// 好友申请操作
	http.HandleFunc("/user/friend.list", api.FriendList)		// 好友列表
	http.HandleFunc("/chatroom/friend_id.get.chatroom_id", api.FriendIdGetChatroomId)		// 通过好友ID 获取 聊天室ID
	http.HandleFunc("/chatroom/get.chatroom.avatar.name.by.chatroom_id", api.GetChatroomAvatarNameByChatRoomID)		// 通过聊天室ID 获取 聊天室头像名称
	http.HandleFunc("/chatroom/chatroom.list", api.ChatroomList)		// 聊天室列表
	http.HandleFunc("/chatroom/create.group", api.CreateGroup)		// 创建群组
	http.HandleFunc("/chatroom/add.group", api.AddGroup)		// 加入群组
	http.HandleFunc("/chatroom/my.group.list", api.MyGroupList)		// 我的群组列表
	http.HandleFunc("/common/get.qiniu.upload.token", api.GetQiniuUploadToken)		// 获取七牛上传token
	http.HandleFunc("/dynamic/publish", api.DynamicPublish)		// 发布动态
	http.HandleFunc("/dynamic/list", api.DynamicList)		// 动态列表
	http.HandleFunc("/app/new.version.get", api.AppNewVersionGet)		// app 最新版本获取

	http.HandleFunc("/common/get.message.id", api.GetMessageId)		// 获取消息ID , 临时使用
	err := http.ListenAndServe(":8066", nil)
	if err != nil {
		panic(err.Error())
	}
}