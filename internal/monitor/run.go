package monitor

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

// websocket

type export struct {
	authorization string
	upgrader      websocket.Upgrader
}

type Content struct {
	WsConn    *websocket.Conn
	closeChan chan byte  // 关闭通知
	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
}

var (
	Export export
)

func Strat() {
	Export = export{
		authorization: "",
		upgrader: websocket.Upgrader{
			//允许跨域
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	// Configure websocket route
	http.HandleFunc("/", HandleConnections)

	// Start the server on localhost port 8000 and log any errors
	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		panic(err)
	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := Export.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Make sure we close the connection when the function returns
	//defer ws.Close()

	//初始化
	var content = Content{
		WsConn:    ws,
		closeChan: make(chan byte),
	}

	//读协程
	go content.wsReadLoop()
}

//读协程 , 处理器
func (content *Content) wsReadLoop() {
	for {
		var receiveStruct interface{}
		// Read in a new message as JSON and map it to a Message object
		err := content.WsConn.ReadJSON(&receiveStruct)

		if err != nil {
			content.wsClose()
			break
		}
		// Send the newly received message to the broadcast channel

	}
}

func (content *Content) wsClose() {
	content.WsConn.Close()
	content.mutex.Lock()
	defer content.mutex.Unlock()
	if !content.isClosed {
		content.isClosed = true
		close(content.closeChan)
	}
}
