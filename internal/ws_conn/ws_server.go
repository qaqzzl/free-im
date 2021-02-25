package ws_conn

import (
	"free-im/pkg/logger"
	"net/http"
)

type wsServer struct {
	Address string // 端口
}

func NewWebSocketServer(address string) *wsServer {
	ws := &wsServer{
		Address: address,
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
