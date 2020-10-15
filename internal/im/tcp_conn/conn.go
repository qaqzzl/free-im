package tcp_conn

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func connRun() {
	addr, err := net.ResolveTCPAddr("tcp", ":1208")
	if err != nil {
		panic(err)
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		print("Fail to start server, %s\n", err)
	}
	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			print("Fail to connect, %s\n", err)
			break
		}
		err = conn.SetKeepAlive(true)
		if err != nil {
			print("Fail to connect, %s\n", err)
			break
		}
		go ConnSocketHandler(conn)
	}
}


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
		WriteChan:	make(chan sendMessage, 1000),
		ReadChan:	make(chan sendMessage, 1000),
	}

	// 消息读取(消息接收) Handler
	go func() {
	Loop:
		for {
			p, err := ctx.Read()
			if err == io.EOF || err != nil {
				ctx.TcpConn.Close()
				break
			}
			// inStr := strings.TrimSpace(string(p.BodyData))
			if p.Action != 99  {
				fmt.Println("接收到的原始消息:", p.Action, p.Version, p.BodyLength, string(p.BodyData))
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
			case ActionAuth: 		// 连接认证

			case ActionMessageRead: // 消息处理
			case ActionClientACK:	// 客户端回执
			case ActionSyncTrigger:	// 消息同步
			case ActionHeadbeat:	// 心跳 5 秒
			case ActionQuit:
				ctx.TcpConn.Close()
				break Loop		// 直接跳出for循环
			default:
				//log.Printf("Unsupported command: %s\n")
			}
		}
	}()
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