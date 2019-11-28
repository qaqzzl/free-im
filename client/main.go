package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"free-im/client/model"
	"net"
	"os"
	"strings"
	"time"
)

func connHandler(c net.Conn) {
	//defer c.Close()

	reader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			input, _ := reader.ReadString('\n')
			inputs := strings.Split(input," ")
			if inputs[0] == "quit" {
				return
			}

			message := model.MessagePackage{}
			switch inputs[0] {
			case "auth":
				content := make(map[string]string)
				content["user_id"] = "1"
				content["access_token"] = "access_token"
				content["device_id"] = uuid.NewV4()
				message = model.MessagePackage {
					Motion:	model.MotionAuth,
					Code: 1,
					Content: content,
				}
			case "pull":

			default:
				content := "啦啦啦 我是普通的文本消息"
				message = model.MessagePackage{
					Motion:	model.MotionMessageSend,
					Code: 1,
					Content: content,
				}
			}

			messagesjon,_ := json.Marshal(message)
			fmt.Println(messagesjon)

			c.Write(messagesjon)
		}
	}()

	go func() {
		for {
			recvData := make([]byte, 2048)
            n, err := c.Read(recvData) //读取数据
            if err != nil {
                fmt.Println(err)
                return
            }
            recvStr := string(recvData[:n])
            fmt.Printf("Response data: %s \n", recvStr)
		}
	}()
}

func main() {
	conn, err := net.Dial("tcp", "localhost:1208")
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
		return
	}

	connHandler(conn)


	<-time.Tick(time.Second * 10000)
}