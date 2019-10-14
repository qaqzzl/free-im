package main

import (
	"encoding/json"
	"fmt"
	"im/server/model"
	"im/server/service"
	"log"
	"net"
	"strings"
)

func connHandler(c net.Conn) {
	if c == nil {
		return
	}

	buf := make([]byte, 4096)

	for {
		cnt, err := c.Read(buf)
		if err != nil || cnt == 0 {
			c.Close()
			break
		}
		inStr := strings.TrimSpace(string(buf[0:cnt]))

		//解析json
		message := db.Message{}
		if err := json.Unmarshal([]byte(inStr), &message); err != nil {
			log.Println(err.Error())
			continue
		}

		ctx := service.Context{
			Conn:    c,
		}

		switch message.Motion {
		case db.MotionAuth:				// 客户端链接认证
			ctx.ClientAuth()
		case db.MotionSendMessage:		// 客户端发送消息
			ctx.ClientSendMessage()
		case db.MotionPullMessage:		// 客户端拉取消息

		case db.MotionQuit:

		default:
			fmt.Printf("Unsupported command: %s\n", message)
		}
	}

	fmt.Printf("Connection from %v closed. \n", c.RemoteAddr())
}

func main() {
	server, err := net.Listen("tcp", ":1208")
	if err != nil {
		fmt.Printf("Fail to start server, %s\n", err)
	}

	fmt.Println("Server Started ...")

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Printf("Fail to connect, %s\n", err)
			break
		}

		go connHandler(conn)
	}
}