package main

import (
	"encoding/json"
	"free-im/library/cache/redis"
	"free-im/library/net/socket"
	"free-im/model"
	"free-im/service"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func connHandler(c net.Conn) {
	if c == nil {
		return
	}
	rconn := redis.GetConn()
	defer rconn.Close()

	buf := make([]byte, 4096)

	ctx := socket.Context{
		ConnSocket:    c,
		InChan:		make(chan *[]byte, 1000),
		OutChan:		make(chan *[]byte, 1000),
	}
	//go func() {
		for {		//读
			cnt, err := c.Read(buf)
			if err != nil || cnt == 0 {
				c.Close()
				break
			}
			inStr := strings.TrimSpace(string(buf[0:cnt]))

			//解析json
			message := model.MessagePackage{}
			if err := json.Unmarshal([]byte(inStr), &message); err != nil {
				log.Println(err.Error())
				continue
			}

			//动作(路由)
			switch message.Action {
			case model.ActionAuth: // 客户端链接认证
				go service.ClientAuth( &ctx )
			case model.ActionMessageSend: // 客户端发送消息
				service.ClientSendMessage(&ctx)
			case model.ActionQuit:
				c.Close()
			default:
				log.Printf("Unsupported command: %s\n")
			}
		}
	//}()

	//消息处理


	//消息发送
	log.Printf("Connection from %v closed. \n", c.RemoteAddr())
}

func main() {
	//http
	go func() {
		http.HandleFunc("/member_id.get.chatroom_id", HttpMemberIdGetChatroomId)

		err := http.ListenAndServe(":8066", nil)
		if ( err != nil ) {
			panic(err.Error())
		}
	}()

	// socket
	server, err := net.Listen("tcp", ":1208")
	if err != nil {
		log.Printf("Fail to start server, %s\n", err)
	}

	log.Println("Server Started ...")

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("Fail to connect, %s\n", err)
			break
		}
		go connHandler(conn)
	}
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

func ErrCheck(err error) {
	if err != nil {
		log.Println(err)
	}
}