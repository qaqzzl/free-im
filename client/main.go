package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"im/client/model"
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

			message := model.Message{}
			switch inputs[0] {
			case "auth":
				message = model.Message{
					Package:     model.MessagePackage{
						Code: 1,
						Content: nil,
					},
					DeviceID:    1,
					UserID:      1,
					Motion:	model.MotionAuth,
				}
			case "pull":

			default:
				message = model.Message{
					Package:     model.MessagePackage{
						Code: 1,
						Content: []byte(input),
					},
					DeviceID:    1,
					UserID:      1,
					Motion:	model.MotionSendMessage,
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