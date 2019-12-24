package main

import (
	"log"
	"net/http"
	"free-im/api/http/v1"
	"os"
)

func main() {
	//http
	//go func() {

	http.HandleFunc("/member_id.get.chatroom_id", v1.MemberIdGetChatroomId)		// 通过会员ID 获取 聊天室ID
	http.HandleFunc("/login", v1.PhoneLogin)		// 手机号登录 / 注册
	http.HandleFunc("/user/member.info", v1.GetUserInfo)		// 获取会员信息
	http.HandleFunc("/user/add.friend", v1.AddFriend)		// 添加好友

		err := http.ListenAndServe(":8066", nil)
		if err != nil {
			panic(err.Error())
		}
	//}()

	// socket
	//server, err := net.Listen("tcp", ":1208")
	//if err != nil {
	//	print("Fail to start server, %s\n", err)
	//}
	//for {
	//	conn, err := server.Accept()
	//	if err != nil {
	//		print("Fail to connect, %s\n", err)
	//		break
	//	}
	//	go connSocketHandler(conn)
	//}
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
