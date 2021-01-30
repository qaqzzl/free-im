package tcp_conn

import (
	"free-im/pkg/logger"
	"net"
)

// TCPServer TCP服务器
type TCPServer struct {
	Address            string // 端口
	AcceptGoroutineNum int    // 接收建立连接的goroutine数量
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(address string, acceptGoroutineNum int) *TCPServer {
	return &TCPServer{
		Address:            address,
		AcceptGoroutineNum: acceptGoroutineNum,
	}
}

// Start 启动服务器
func (t *TCPServer) Start() {
	addr, err := net.ResolveTCPAddr("tcp", t.Address)
	if err != nil {
		panic(err)
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		print("Fail to start server, %s\n", err)
	}
	for i := 0; i < t.AcceptGoroutineNum; i++ {
		go t.Accept(server)
	}
	logger.Sugar.Info("tcp server start")
	select {}
}

// Accept 接收客户端的TCP长连接的建立
func (t *TCPServer) Accept(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			print("Fail to connect, %s\n", err)
			break
		}
		err = conn.SetKeepAlive(true)
		if err != nil {
			print("Fail to connect, %s\n", err)
			break
		}
		connContext := NewConnContext(conn)
		go connContext.DoConn()
	}
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
