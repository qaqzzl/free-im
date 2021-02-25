package ws_conn

import "net/http"

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
	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		panic(err)
	}
}
