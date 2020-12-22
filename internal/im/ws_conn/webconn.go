package ws_conn

import (
	"free-im/pkg/logger"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Content struct {
	RedisConn   redis.Conn
	WsConn      *websocket.Conn
	WsConnIndex string     // 房间index
	closeChan   chan byte  // 关闭通知
	mutex       sync.Mutex // 避免重复关闭管道
	isClosed    bool
	//inChan chan *[]byte								// 读队列
	outChan chan *[]byte // 写队列
}

// Configure the upgrader
var (
	//完成握手操作
	upgrader = websocket.Upgrader{
		//允许跨域(一般来讲,websocket都是独立部署的)
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
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
func (content *Content) wsReadLoop() {
	defer content.RedisConn.Close() //处理器退出后关闭redis连接资源
	for {
		var receiveStruct interface{}
		// Read in a new message as JSON and map it to a Message object
		err := content.WsConn.ReadJSON(&receiveStruct)

		if err != nil {
			content.wsClose()
			break
		}
		// Send the newly received message to the broadcast channel
		switch receiveStruct {
		case "ping":
			// log.Println(receiveStruct.ClientType, receiveStruct.UserId)
			msg := []byte(`{"code":200,"msg":"服务器响应正常"}`)
			content.outChan <- &msg
		default:
			logger.Sugar.Error("错误的请求方法")
		}
	}
}

//写协程
func (content *Content) wsWriteLoop() {
	for {
		select {
		// 取一个应答
		case msg := <-content.outChan:
			// 写给websocket
			if err := content.WsConn.WriteMessage(1, *msg); err != nil {
				goto error
			}
		case <-content.closeChan:
			goto closed
		}
	}
error:
	content.WsConn.Close()
closed:
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
