package connect

import (
	"encoding/json"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"github.com/gorilla/websocket"
	"net/http"
)

type wsServer struct {
	Address string // 端口
}

func NewWebSocketServer(address string) *wsServer {
	WSServer = &wsServer{
		Address: address,
	}
	return WSServer
}

var WSServer *wsServer

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

func (ws *wsServer) Start() {
	logger.Sugar.Info("ws server start", ws.Address)

	// Configure websocket route
	http.HandleFunc("/", Connections)

	// Start the server on localhost port 8000 and log any errors
	err := http.ListenAndServe(ws.Address, nil)
	if err != nil {
		panic(err)
	}
}

func Connections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	//初始化
	var conn = &Conn{
		CoonType: ConnTypeWS,
		WS:       ws,
	}

	//读协程
	go wsReadLoop(conn)
}

// 读协程 处理WS连接
func wsReadLoop(conn *Conn) {
	// Make sure we close the connection when the function returns
	defer conn.WS.Close()
	for {
		var mp pbs.MsgPackage
		// Read in a new message as JSON and map it to a Message object
		// err := conn.c.ReadJSON(mp)
		_, data, err := conn.WS.ReadMessage()
		logger.Sugar.Debug(string(data))
		if err != nil {
			conn.Close()
			break
		}
		if err = json.Unmarshal(data, &mp); err != nil {
			logger.Sugar.Error(err)
			continue
		}
		// Send the newly received message to the broadcast channel
		Handler.Handler(conn, mp)
	}
}
