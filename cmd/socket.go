package main

import (
	"encoding/json"
	"free-im/library/cache/redis"
	"free-im/library/net/socket"
	"free-im/model"
	"free-im/service"
	"log"
	"net"
	"strings"
)

func connSocketHandler(c net.Conn) {
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
