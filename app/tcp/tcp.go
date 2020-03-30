package tcp

import (
	"encoding/json"
	"fmt"
	"free-im/model"
	"free-im/service/tcp"
	"log"
	"net"
	"strings"
)


func ConnSocketHandler(c net.Conn) {
	if c == nil {
		return
	}

	buf := make([]byte, 4096)

	ctx := tcp.Context{
		TcpConn:    c,
		InChan:		make(chan *[]byte, 1000),
		OutChan:	make(chan *[]byte, 1000),
	}
	for {		//读
		cnt, err := c.Read(buf)
		if err != nil || cnt == 0 {
			c.Close()
			break
		}
		inStr := strings.TrimSpace(string(buf[0:cnt]))
		fmt.Println("接收到的原始消息:",inStr)
		//动作(路由)
		switch inStr[0:2] {
		case model.ActionAuth: // 客户端链接认证
			fmt.Println("接收到的原始消息:",cnt)
			//解析json
			message := model.AuthMessage{}
			if err := json.Unmarshal([]byte(inStr[2:cnt]), &message); err != nil {
				log.Println(err.Error())
				continue
			}
			go tcp.ClientAuth( &ctx, message )
		case model.ActionMessageSend: // 客户端发送消息
			//解析json
			message := model.MessagePackage{}
			if err := json.Unmarshal([]byte(inStr[2:cnt]), &message); err != nil {
				log.Println(err.Error())
				continue
			}
			tcp.ClientSendMessage(&ctx, message)
		case model.ActionQuit:
			c.Close()
		default:
			//log.Printf("Unsupported command: %s\n")
		}
	}

	//消息处理


	//消息发送


	// 连接断开处理
	log.Printf("Connection from %v closed. \n", c.RemoteAddr())
}