package main

import (
	"flag"
	"fmt"
	httpV1 "free-im/api/http_conn/v1"
	"free-im/configs"
	"free-im/internal/im/ws_conn"
	"log"
	"net/http"
	"os"
)



func main() {
	// 配置文件
	ConfPath := flag.String("cpath", "./config.conf", "config file")
	config.InitConfig(*ConfPath)

	// http
	go init_http()
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
	http.HandleFunc("/login", httpV1.PhoneLogin)		// 手机号登录 / 注册
	http.HandleFunc("/login/send.login.sms", httpV1.SendLoginSms)		// 发送登录手机号验证码
	http.HandleFunc("/user/member.info", httpV1.GetMemberInfo)		// 获取会员信息
	http.HandleFunc("/user/update.member.info", httpV1.UpdateMemberInfo)		// 修改会员信息
	http.HandleFunc("/user/others.home.info", httpV1.OthersHomeInfo)		// 获取用户基本信息(他人主页)
	http.HandleFunc("/search/friend", httpV1.SearchMember)		// 搜索好友
	http.HandleFunc("/user/add.friend", httpV1.AddFriend)		// 添加好友
	http.HandleFunc("/user/del.friend", httpV1.DelFriend)		// 删除好友
	http.HandleFunc("/user/friend.apply.list", httpV1.FriendApplyList)		// 好友申请列表
	http.HandleFunc("/user/friend.apply.action", httpV1.FriendApplyAction)		// 好友申请操作
	http.HandleFunc("/user/friend.list", httpV1.FriendList)		// 好友列表
	http.HandleFunc("/chatroom/friend_id.get.chatroom_id", httpV1.FriendIdGetChatroomId)		// 通过好友ID 获取 聊天室ID
	http.HandleFunc("/chatroom/get.chatroom.avatar.name.by.chatroom_id", httpV1.GetChatroomAvatarNameByChatRoomID)		// 通过聊天室ID 获取 聊天室头像名称
	http.HandleFunc("/chatroom/chatroom.list", httpV1.ChatroomList)		// 聊天室列表
	http.HandleFunc("/chatroom/create.group", httpV1.CreateGroup)		// 创建群组
	http.HandleFunc("/chatroom/add.group", httpV1.AddGroup)		// 加入群组
	http.HandleFunc("/chatroom/my.group.list", httpV1.MyGroupList)		// 我的群组列表
	http.HandleFunc("/common/get.qiniu.upload.token", httpV1.GetQiniuUploadToken)		// 获取七牛上传token
	http.HandleFunc("/dynamic/publish", httpV1.DynamicPublish)		// 发布动态
	http.HandleFunc("/dynamic/list", httpV1.DynamicList)		// 动态列表
	http.HandleFunc("/app/new.version.get", httpV1.AppNewVersionGet)		// app 最新版本获取

	http.HandleFunc("/common/get.message.id", httpV1.GetMessageId)		// 获取消息ID , 临时使用
	err := http.ListenAndServe(":8066", nil)
	if err != nil {
		panic(err.Error())
	}
}

func init_im_tcp() {

}

func init_im_ws() {
	// Configure websocket route
	http.HandleFunc("/", ws_conn.HandleConnections)

	// Start listening for incoming chat messages
	//go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ")
		panic( err )
	}
}