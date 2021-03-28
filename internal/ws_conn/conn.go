package ws_conn

import (
	"free-im/pkg/protos/pbs"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Conn struct {
	c         *websocket.Conn
	closeChan chan byte  // 关闭通知
	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	//inChan chan *[]byte								// 读队列
	outChan    chan *[]byte // 写队列
	Version    int32
	DeviceID   string // 设备id 简写 DID
	UserID     string // 用户id 简写 UID
	DeviceType string // 设备类型, 移动端:mobile , PC端:pc
	ClientType string // 客户端类型, (android, ios) | (windows, mac, linux)
	IsAuth     bool   // 是否认证(登录)
}

// Configure the upgrader
var (
	//完成握手操作
	upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Connections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Make sure we close the connection when the function returns
	//defer ws.Close()

	//初始化
	var conn = Conn{
		c:         ws,
		outChan:   make(chan *[]byte, 100),
		closeChan: make(chan byte),
	}

	//读协程
	go conn.wsReadLoop()

	// 写协程
	go conn.wsWriteLoop()
}

//读协程 , 处理器
func (conn *Conn) wsReadLoop() {
	for {
		var mp pbs.MsgPackage
		// Read in a new message as JSON and map it to a Message object
		err := conn.c.ReadJSON(&mp)

		if err != nil {
			conn.Close()
			break
		}

		// Send the newly received message to the broadcast channel
		Handler.Handler(conn, mp)
	}
}

//写协程
func (conn *Conn) wsWriteLoop() {
	for {
		select {
		// 取一个应答
		case msg := <-conn.outChan:
			// 写给websocket
			if err := conn.c.WriteMessage(1, *msg); err != nil {
				goto error
			}
		case <-conn.closeChan:
			goto closed
		}
	}
error:
	conn.c.Close()
closed:
}

func (conn *Conn) Close() {
	//conn.WsConn.Close()
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if !conn.isClosed {
		conn.isClosed = true
		close(conn.closeChan)
	}
}
