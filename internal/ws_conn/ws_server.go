package ws_conn

import (
	"free-im/pkg/logger"
	cmap "github.com/orcaman/concurrent-map"
	"net/http"
)

type wsServer struct {
	Address        string             // 端口
	ServerConnPool cmap.ConcurrentMap // 链接池
}

func NewWebSocketServer(address string) *wsServer {
	ws := &wsServer{
		Address:        address,
		ServerConnPool: cmap.New(),
	}
	return ws
}

func (ws *wsServer) Start() {
	// Configure websocket route
	http.HandleFunc("/", Connections)

	// Start listening for incoming chat messages
	//go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	logger.Sugar.Info("ws server start")
	err := http.ListenAndServe(ws.Address, nil)
	if err != nil {
		panic(err)
	}
}

// load 获取链接
func (ws *wsServer) LoadConn(UserID string, DeviceID string) (ctx *Context) {
	tmp, ok := ws.ServerConnPool.Get(UserID)
	if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
		for _, vo := range tmp.(cmap.ConcurrentMap).Items() {
			ctx := vo.(*Context)
			if ctx.DeviceID == DeviceID {
				break
			}
		}
	}
	return ctx
}

func (ws *wsServer) LoadConnsByUID(UserID string) (ctxs []*Context) {
	tmp, ok := ws.ServerConnPool.Get(UserID)
	if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
		for _, vo := range tmp.(cmap.ConcurrentMap).Items() {
			ctx := vo.(*Context)
			ctxs = append(ctxs, ctx)
		}
	}
	return ctxs
}

// store 存储
func (ws *wsServer) StoreConn(ctx *Context) {
	if tmp, ok := ws.ServerConnPool.Get(ctx.UserID); ok {
		device_map := tmp.(cmap.ConcurrentMap)
		device_map.Set(ctx.DeviceType, ctx)
		ws.ServerConnPool.Set(ctx.UserID, device_map)
	} else {
		device_map := cmap.New()
		device_map.Set(ctx.DeviceType, ctx)
		ws.ServerConnPool.Set(ctx.UserID, device_map)
	}
}
