package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"free-im/server/model"
	uuid "github.com/satori/go.uuid"
	"net"
	"os"
	"strings"
)

func connHandler(c net.Conn) {
	//defer c.Close()

	reader := bufio.NewReader(os.Stdin)

	//聊天室ID
	chatroom_id := ""
	go func() {		//发送消息
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
				content["device_id"] = uuid.NewV4().String()
				message = model.MessagePackage {
					Motion:	model.MotionAuth,
					Code: 1,
					Content: content,
				}
			case "pull":

			default:
				content := make(map[string]string)
				content["user_id"] = "1"
				content["chatroom_id"] = chatroom_id
				content["device_id"] = uuid.NewV4().String()
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

	go func() {		//接收消息
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