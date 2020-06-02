package tcp_conn

import (
	"bufio"
	"encoding/json"
	"free-im/dao"
	"github.com/orcaman/concurrent-map"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)


func ConnSocketHandler(conn net.Conn) {
	// fmt.Println("接入新连接")
	if conn == nil {
		return
	}

	// 心跳时间
	headbeatTime := time.Now().Unix()

	ctx := Context{
		TcpConn:    conn,
		r:		bufio.NewReader(conn),
		IsConnStatus:    true,
		WriteChan:	make(chan sendMessage, 1000),
		ReadChan:	make(chan sendMessage, 1000),
		Status: true,
	}
	logic := Logic{
		ctx: &ctx,
	}

	//defer redisconn.Close()
	//ack超时 消息重传
	go func() {
		redisconn := dao.NewRedis()
		ticker := time.NewTicker(time.Second * 1)
		for {
			if ctx.IsConnStatus == false {
				break
			}
			if ctx.IsAuth == false {
				<-ticker.C
				continue
			}
			if message_id,err := redisconn.Do("RPOP","list_message_ack_timeout_retransmit:"+ctx.UserID+":"+ctx.ClientType); err != nil {
				log.Println("list_message_ack_timeout_retransmit err", err.Error())
				<-ticker.C
				continue
			} else {
				if message_id == nil {
					<-ticker.C
					continue
				}
				if message, err := redisconn.Do("HGET","hash_message_ack_timeout_retransmit:"+ctx.UserID+":"+ctx.ClientType, message_id); err != nil {
					log.Println("hash_message_ack_timeout_retransmit err", err.Error())
					<-ticker.C
					continue
				} else {
					if message == nil {
						continue
					}
					redisconn.Do("RPUSH", "list_message_ack_timeout_retransmit:"+ctx.UserID+":"+ctx.ClientType, message_id)
					messages := strings.Split(string(message.([]byte)),"|")
					timeUnix, _ := strconv.ParseInt(messages[0], 10, 64)
					UnixNano := time.Now().UnixNano() / 1e6
					if timeUnix >  UnixNano {
						continue
					}
					redisconn.Do("HSET", "hash_message_ack_timeout_retransmit:"+ctx.UserID+":"+ctx.ClientType, message_id, strconv.FormatInt(UnixNano + 1000,10)+"|"+messages[1])
					logic.ctx.WriteChan <- sendMessage{
						Conn:ctx.TcpConn,
						Package:Package{
							Version: Version,
							Action: ActionMessageSend,
							BodyData: []byte(messages[1]),
						},
					}
				}
			}
		}
	}()

	// 离线消息
	go func() {
		redisconn := dao.NewRedis()
		ticker := time.NewTicker(time.Second * 1)
		for {
			if ctx.IsConnStatus == false {
				break
			}
			if ctx.IsAuth == false {
				<-ticker.C
				continue
			}
			if message,err := redisconn.Do("RPOP","list_message_offline:"+ctx.UserID); err != nil {
				log.Println("list_message_offline err", err.Error())
				<-ticker.C
				continue
			} else {
				if message == nil {
					<-ticker.C
					continue
				}
				packages := Package{
					Version: Version,
					Action: ActionMessageSend,
					BodyData: message.([]byte),
				}
				m := MessagePackage{}
				if err := json.Unmarshal(packages.BodyData, &m); err != nil {
					log.Println(err.Error())
					continue
				}
				tmp, ok := SocketConnPool.Get(ctx.UserID)
				if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
					for  _, vo := range tmp.(cmap.ConcurrentMap).Items() {
						timeUnix := time.Now().UnixNano()  / 1e6
						device := vo.(ClientDevice)
						redisconn.Do("LPUSH", "list_message_ack_timeout_retransmit:"+ctx.UserID+":"+device.ClientType, m.MessageId)
						redisconn.Do("HSET", "hash_message_ack_timeout_retransmit:"+ctx.UserID+":"+device.ClientType, m.MessageId, strconv.FormatInt(timeUnix + 1000,10)+"|"+string(packages.BodyData))
						logic.ctx.WriteChan <- sendMessage{
							Conn:device.Conn,
							Package:packages,
						}
					}
				} else {
					// fmt.Println("设备未在线 , 未读消息写入redis")
					redisconn.Do("LPUSH", "list_message_offline:"+ctx.UserID, packages.BodyData)
				}
			}
		}
	}()

	//消息发送 Handler
	go func() {
		for v := range ctx.WriteChan {
			if ctx.IsConnStatus == false {
				break
			}
			logic.SendResponse(v.Conn, v.Package)
		}
	}()

	// 消息读取(消息接收) Handler
	go func() {
	Loop:
		for {
			p, err := ctx.Read()
			if err == io.EOF || err != nil {
				logic.Close()
				break
			}
			// inStr := strings.TrimSpace(string(p.BodyData))
			if p.Action != 99  {
				// fmt.Println("接收到的原始消息:", p.Action, p.Version, p.SequenceId, string(p.BodyData))
			}
			headbeatTime = time.Now().Unix()
			// 动作(路由)
			switch p.Action {
			case ActionMessageID:	// 获取消息ID
				//m := MessagePackage{}
				//if err := json.Unmarshal(p.BodyData, &m); err != nil {
				//	log.Println(err.Error())
				//	continue
				//}
				//logic.GetMessageId(m)
			case ActionAuth: // 客户端链接认证
				m := AuthMessage{}
				if err := json.Unmarshal(p.BodyData, &m); err != nil {
					log.Println(err.Error())
					continue
				}
				logic.ClientAuth( m )
			case ActionMessageRead: // 客户端发送消息
				m := MessagePackage{}
				if err := json.Unmarshal(p.BodyData, &m); err != nil {
					log.Println(err.Error())
					continue
				}
				logic.ClientMessage(m)

			case ActionClientACK:	// 客户端回执
				logic.ClientACK(p)
			case ActionSyncTrigger:	// 消息同步

			case ActionHeadbeat:	// 心跳 5 秒
				Package := Package{
					Version: Version,
					Action: ActionHeadbeatACK,
					BodyData: []byte(""),
				}
				logic.SendResponse(ctx.TcpConn, Package)
			case ActionQuit:
				logic.Close()
				break Loop		// 直接跳出for循环
			default:
				//log.Printf("Unsupported command: %s\n")
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second * 1)
		for  {
			if (headbeatTime+10) < time.Now().Unix() {
				logic.Close()
				break
			}
			<- ticker.C
		}
	}()

	ticker := time.NewTicker(time.Second * 1)
	for {
		// 如果连接关闭并且OutChan为空, 则退出
		if len(ctx.WriteChan) == 0 && ctx.IsConnStatus == false {
			//fmt.Println("断开处理")
			//fmt.Printf("Connection from %v closed. \n", conn.RemoteAddr())
			logic.Close()
			break
		} else {
			<-ticker.C
		}
	}

	ctx.Status = false
}


// 系统监听
func SystemMonitor() {
	//go func() {
	//	ticker := time.NewTicker(time.Second * 3)
	//	for {
	//		<-ticker.C
	//		fmt.Println("-----------------------------------")
	//		fmt.Println("连接用户数: ",SocketConnPool.Count())
	//		for key,vo := range SocketConnPool.Items() {
	//			fmt.Println("--------------")
	//			fmt.Println("连接用户ID: ", key)
	//			ConcurrentMap := vo.(cmap.ConcurrentMap)
	//			for k,v := range ConcurrentMap.Items() {
	//				fmt.Println("连接设备类型: ", k)
	//				fmt.Println("连接设备ID: ", v.(ClientDevice).DeviceID)
	//			}
	//
	//		}
	//	}
	//}()
}