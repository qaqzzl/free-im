package main

import (
	"encoding/json"
	"fmt"
	"free-im/server/library/cache/redis"
	"free-im/server/model"
	"free-im/server/service"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func connHandler(c net.Conn) {
	if c == nil {
		return
	}

	buf := make([]byte, 4096)

	ctx := model.Context{
		ConnSocket:    c,
		InChan:		make(chan *model.MessagePackage, 1000),
		OutChan:		make(chan *model.MessagePackage, 1000),
	}
	//go func() {
		for {		//读
			cnt, err := c.Read(buf)
			if err != nil || cnt == 0 {
				c.Close()
				break
			}
			inStr := strings.TrimSpace(string(buf[0:cnt]))

			//解析json
			if err := json.Unmarshal([]byte(inStr), &ctx.Message); err != nil {
				log.Println(err.Error())
				continue
			}

			//动作(路由)
			switch ctx.Message.Motion {
			case model.MotionAuth: // 客户端链接认证
				go service.ClientAuth( &ctx )
			case model.MotionMessageSend: // 客户端发送消息
				service.ClientSendMessage(&ctx)
			case model.MotionQuit:
				c.Close()
			default:
				log.Printf("Unsupported command: %s\n")
			}
		}
	//}()

	//消息处理


	//消息发送
	log.Printf("Connection from %v closed. \n", c.RemoteAddr())
}

func main() {
	//http
	go func() {
		http.HandleFunc("/member_id.get.chatroom_id", HttpMemberIdGetChatroomId)
		err := http.ListenAndServe(":8066", nil)
		if ( err != nil ) {
			panic(err.Error())
		}
	}()

	// socket
	server, err := net.Listen("tcp", ":1208")
	if err != nil {
		log.Printf("Fail to start server, %s\n", err)
	}

	log.Println("Server Started ...")

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("Fail to connect, %s\n", err)
			break
		}
		go connHandler(conn)
	}
}

// HTTP 通过会员ID 获取 聊天室ID
func HttpMemberIdGetChatroomId(writer http.ResponseWriter, request *http.Request) {
	rconn := redis.GetConn()
	var err error

	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	var field string
	if formData["user_id"].(string) > formData["member_id"].(string) {
		field = formData["user_id"].(string) +","+ formData["member_id"].(string)
	} else {
		field = formData["member_id"].(string) +","+ formData["user_id"].(string)
	}
	var res interface{}
	if res,err = rconn.Do("HGET", "hash_im_chatroom_member_id_get_chatroom_id", field); err == nil {
		log.Println(err)
	}

	var chatroom_id string
	if res == nil {
		//生成聊天室ID
		chatroom_id = uuid.NewV4().String()
		rconn.Do("SADD", "set_im_chatroom_member_"+chatroom_id, formData["user_id"], formData["member_id"])			//创建聊天室
		rconn.Do("HSET", "hash_im_chatroom_member_id_get_chatroom_id", field, chatroom_id)			//创建聊天室
	} else {
		chatroom_id = string( res.([]uint8) )
	}

	requestBody := fmt.Sprintf(`{
"chatroom_id":"%s",
"status": "%s",
"code": %d
}`,
		chatroom_id, "ok",0)
	rconn.Close()
	writer.Write([]byte(requestBody))
}

// HTTP 获取好友列表

// HTTP 获取聊天室列表

// HTTP 获取群聊列表

// HTTP 添加好友

// HTTP 删除好友

func init() {
	file := "./" +"message"+ ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func ErrCheck(err error) {
	if err != nil {
		log.Println(err)
	}
}