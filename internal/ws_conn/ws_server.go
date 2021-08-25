package ws_conn

import (
	"free-im/pkg/logger"
	"free-im/pkg/service/user"
	cmap "github.com/orcaman/concurrent-map"
	"net/http"
	"strconv"
)

type wsServer struct {
	Address        string             // 端口
	ServerConnPool cmap.ConcurrentMap // 链接池
}

var WSServer *wsServer

func NewWebSocketServer(address string) *wsServer {
	WSServer = &wsServer{
		Address:        address,
		ServerConnPool: cmap.New(),
	}
	return WSServer
}

func (ws *wsServer) Start() {
	// Configure websocket route
	http.HandleFunc("/", Connections)

	// Start listening for incoming chat messages
	//go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	logger.Sugar.Info("ws server start", ws.Address)
	err := http.ListenAndServe(ws.Address, nil)
	if err != nil {
		panic(err)
	}
}

// load 获取链接
func (ws *wsServer) LoadConn(UserID string, DeviceID string) (conn *Conn) {
	tmp, ok := ws.ServerConnPool.Get(UserID)
	if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
		for _, vo := range tmp.(cmap.ConcurrentMap).Items() {
			conn = vo.(*Conn)
			if conn.DeviceID == DeviceID {
				break
			}
		}
	}
	return conn
}

func (ws *wsServer) LoadConnsByUID(UserID int64) (conns []*Conn) {
	tmp, ok := ws.ServerConnPool.Get(strconv.Itoa(int(UserID)))
	if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
		for _, vo := range tmp.(cmap.ConcurrentMap).Items() {
			conn := vo.(*Conn)
			conns = append(conns, conn)
		}
	}
	return conns
}

// store 存储
func (ws *wsServer) StoreConn(conn *Conn) {
	key := strconv.Itoa(int(conn.UserID))
	if tmp, ok := ws.ServerConnPool.Get(key); ok {
		device_map := tmp.(cmap.ConcurrentMap)
		device_map.Set(conn.DeviceType, conn)
		ws.ServerConnPool.Set(key, device_map)
	} else {
		device_map := cmap.New()
		device_map.Set(conn.DeviceType, conn)
		ws.ServerConnPool.Set(key, device_map)
	}
	// 用户在线状态
	user.User.SetUserOnline(conn.UserID, true, conn.DeviceType)
}
