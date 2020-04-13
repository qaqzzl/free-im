package socket

import (
	"bufio"
	"encoding/json"
	"fmt"
	"free-im/app/model"
	"free-im/dao"
	"io"
	"log"
	"net"
	"strings"
	"time"
)


func ConnSocketHandler(conn net.Conn) {
	if conn == nil {
		return
	}

	ctx := Context{
		TcpConn:    conn,
		r:		bufio.NewReader(conn),
		IsConnStatus:    true,
		InChan:		make(chan sendMessage, 1000),
		OutChan:	make(chan sendMessage, 1000),
	}
	message := Message{
		ctx: &ctx,
	}
	redisconn := dao.NewRedis()
	//defer redisconn.Close()
	//消息处理 Handler
	go func() {
		ticker := time.NewTicker(time.Second * 2)
		for {
			if ctx.IsConnStatus == false {
				break
			}
			if ctx.IsAuth == false {
				continue
			}
			if value,err := redisconn.Do("RPOP","list_message_send_failure:"+ctx.UserID); err != nil {
				fmt.Println("list_message_send_failure err", err.Error())
				continue
			} else {
				if value != nil {
					message.ctx.OutChan <- sendMessage{
						Conn:ctx.TcpConn,
						ReceiveUserID:ctx.UserID,
						Message:value.([]byte),
					}
				} else {
					<-ticker.C
				}
			}
		}
	}()


	//消息发送 Handler
	go func() {
		for v := range ctx.OutChan {
			if ctx.IsConnStatus == false {
				// 消息发送失败 , 记录下
				redisconn.Do("LPUSH", "list_message_send_failure:"+v.ReceiveUserID, v.Message)
			} else {
				message.SendResponse(v.Conn, v.ReceiveUserID, v.Message)
			}
		}
	}()

	// 消息读取(消息接收) Handler
	//buf := make([]byte, 4096)
	go func() {
	Loop:
		for {
			buf, err := ctx.Read()
			if err == io.EOF{
				TcpClose(&ctx)
				break
			}
			if err != nil {
				continue
			}
			inStr := strings.TrimSpace(string(buf))
			fmt.Println("接收到的原始消息:",inStr)
			fmt.Println("UserID:",ctx.UserID)
			//动作(路由)
			switch inStr[0:2] {
			case model.ActionAuth: // 客户端链接认证
				m := model.AuthMessage{}
				if err := json.Unmarshal([]byte(inStr[2:]), &m); err != nil {
					log.Println(err.Error())
					continue
				}
				message.ClientAuth( m )
			case model.ActionMessage: // 客户端发送消息
				//解析json
				m := model.MessagePackage{}
				if err := json.Unmarshal([]byte(inStr[2:]), &m); err != nil {
					log.Println(err.Error())
					continue
				}
				message.ClientMessage(m)
			case model.ActionSyncTrigger:	// 消息同步

			case model.ActionQuit:
				message.Close()
				break Loop		// 直接跳出for循环
			default:
				//log.Printf("Unsupported command: %s\n")
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 1)
	for {
		// 如果连接关闭并且OutChan为空, 则退出
		if len(ctx.OutChan) == 0 && ctx.IsConnStatus == false {
			fmt.Println("断开处理")
			fmt.Printf("Connection from %v closed. \n", conn.RemoteAddr())
			message.Close()
			break
		} else {
			<-ticker.C
		}
	}
}


// 系统监听
func SystemMonitor() {
	//go func() {
	//	ticker := time.NewTicker(time.Second * 3)
	//	for {
	//		<-ticker.C
	//		fmt.Println("连接用户数: ",len(model.SocketConnPool))
	//		for key,vo := range model.SocketConnPool {
	//			fmt.Println("-----------------------------------")
	//			fmt.Println("连接用户ID: ", key)
	//			for k,_ := range vo {
	//				fmt.Println("连接设备DeviceID: ", k)
	//			}
	//
	//		}
	//	}
	//}()
}