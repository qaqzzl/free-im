package tcp

import (
	"encoding/json"
	"fmt"
	"free-im/app/model"
	"free-im/app/service/tcp"
	"log"
	"net"
	"strings"
	"time"
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

	//消息处理 Handler


	//消息发送 Handler


	for {		//消息读取(消息接收) Handler
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

	// 连接断开处理
	log.Printf("Connection from %v closed. \n", c.RemoteAddr())
}


// 系统监听
func SystemMonitor() {
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for {
			<-ticker.C
			fmt.Println("连接用户数: ",len(model.SocketConnPool))
			for key,vo := range model.SocketConnPool {
				fmt.Println("-----------------------------------")
				fmt.Println("连接用户ID: ", key)
				for k,_ := range vo {
					fmt.Println("连接设备DeviceID: ", k)
				}

			}
		}
	}()
}