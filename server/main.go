package server

import (
	"encoding/json"
	"free-im/server/model"
	"free-im/server/service"
	"log"
	"net"
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
	}
	for {
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

		//路由
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

	log.Printf("Connection from %v closed. \n", c.RemoteAddr())
}

func main() {

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