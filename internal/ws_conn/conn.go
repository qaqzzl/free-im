package ws_conn

import (
	"encoding/json"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/service/user"
	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Conn struct {
	c          *websocket.Conn
	closeChan  chan byte  // 关闭通知
	mutex      sync.Mutex // 避免重复关闭管道
	isClosed   bool
	Version    int32
	DeviceID   string // 设备id 简写 DID
	UserID     int64  // 用户id 简写 UID
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
	log.Println("new conn")
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

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
	// Make sure we close the connection when the function returns
	defer conn.c.Close()
	for {
		var mp pbs.MsgPackage
		// Read in a new message as JSON and map it to a Message object
		// err := conn.c.ReadJSON(mp)
		_, data, err := conn.c.ReadMessage()
		logger.Sugar.Info(string(data))
		if err != nil {
			conn.Close()
			break
		}
		if err = json.Unmarshal(data, &mp); err != nil {
			logger.Sugar.Info(err)
			continue
		}
		// Send the newly received message to the broadcast channel
		Handler.Handler(conn, mp)
	}
}

func (conn *Conn) Close() {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if !conn.isClosed {
		if conn.IsAuth {
			key := strconv.Itoa(int(conn.UserID))
			if user_map, ok := WSServer.ServerConnPool.Get(key); ok {
				user_map.(cmap.ConcurrentMap).Remove(conn.DeviceType)
				WSServer.ServerConnPool.Set(key, user_map)
			}
			// 用户在线状态
			user.User.SetUserOnline(conn.UserID, false, conn.DeviceType)
		}
		conn.isClosed = true
		close(conn.closeChan)
		conn.c.Close()
	}
}

func (conn *Conn) Write(mp pbs.MsgPackage) error {
	err := conn.c.WriteJSON(mp)
	return err
}
