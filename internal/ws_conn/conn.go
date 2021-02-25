package ws_conn

import (
	"free-im/pkg/logger"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Content struct {
	WsConn    *websocket.Conn
	closeChan chan byte  // 关闭通知
	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	//inChan chan *[]byte								// 读队列
	outChan chan *[]byte // 写队列
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
	var content = Content{
		WsConn:    ws,
		outChan:   make(chan *[]byte, 100),
		closeChan: make(chan byte),
	}

	//读协程
	go content.wsReadLoop()

	// 写协程
	go content.wsWriteLoop()
}

//读协程 , 处理器
func (ctx *Content) wsReadLoop() {
	for {
		var receiveStruct interface{}
		// Read in a new message as JSON and map it to a Message object
		err := ctx.WsConn.ReadJSON(&receiveStruct)

		if err != nil {
			ctx.wsClose()
			break
		}
		// Send the newly received message to the broadcast channel
		switch receiveStruct {
		case "ping":
			// log.Println(receiveStruct.ClientType, receiveStruct.UserId)
			msg := []byte(`{"code":200,"msg":"服务器响应正常"}`)
			ctx.outChan <- &msg
		default:
			logger.Sugar.Error("错误的请求方法")
		}
	}
}

//写协程
func (ctx *Content) wsWriteLoop() {
	for {
		select {
		// 取一个应答
		case msg := <-ctx.outChan:
			// 写给websocket
			if err := ctx.WsConn.WriteMessage(1, *msg); err != nil {
				goto error
			}
		case <-ctx.closeChan:
			goto closed
		}
	}
error:
	ctx.WsConn.Close()
closed:
}

func (ctx *Content) wsClose() {
	//ctx.WsConn.Close()
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if !ctx.isClosed {
		ctx.isClosed = true
		close(ctx.closeChan)
	}
}
