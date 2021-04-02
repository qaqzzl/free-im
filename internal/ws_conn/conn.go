package ws_conn

import (
	"free-im/pkg/protos/pbs"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Conn struct {
	c          *websocket.Conn
	closeChan  chan byte  // 关闭通知
	mutex      sync.Mutex // 避免重复关闭管道
	isClosed   bool
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
		closeChan: make(chan byte),
	}

	//读协程
	go conn.wsReadLoop()
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

func (conn *Conn) Close() {
	//conn.WsConn.Close()
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if !conn.isClosed {
		conn.isClosed = true
		close(conn.closeChan)
	}
}

func (conn *Conn) Write(mp pbs.MsgPackage) error {
	err := conn.c.WriteJSON(mp)
	return err
}
