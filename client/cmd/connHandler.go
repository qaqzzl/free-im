package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"free-im/client/model"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

func connHandler(c net.Conn) {
	//defer c.Close()

	reader := bufio.NewReader(os.Stdin)

	//聊天室ID
	chatroom_id := ""
	is_auth := false
	go func() {		//发送消息
		for {
			input, _ := reader.ReadString('\n')
			input = string(bytes.TrimRight([]byte(input), "\n"))
			//inputs := strings.Split(input," ")
			//if inputs[0] == "quit" {
			//	return
			//}
			message := model.MessagePackage{}

			if chatroom_id == "" {
				requestBody := fmt.Sprintf(`{
"user_id":"%s",
"access_token": "%s",
"member_id": "%s"
}`,
input, "access_token", "1")
				fmt.Println(requestBody)
				var jsonStr = []byte(requestBody)
				url := "http://127.0.0.1:8066/member_id.get.chatroom_id"
				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				// 初始化请求变量结构
				formData := make(map[string]interface{})
				// 调用json包的解析，解析请求body
				json.Unmarshal(body, &formData)
				chatroom_id = formData["chatroom_id"].(string)		//聊天室ID
			}

			if is_auth == false {
				content := make(map[string]string)
				content["user_id"] = "1"
				content["access_token"] = "access_token"
				content["device_id"] = uuid.NewV4().String()
				message := model.MessagePackage {
					Motion:	10,
					Content: content,
				}
				messagesjon,_ := json.Marshal(message)
				c.Write(messagesjon)
				<- time.Tick(time.Second * 1)
				is_auth = true
			}



			<- time.Tick(time.Millisecond * 300)
			content := "啦啦啦  我是消息内容 . " + uuid.NewV4().String()
			message = model.MessagePackage{
				Motion:	6,
				ChatroomId: chatroom_id,
				ClassCode: 1,
				Content: content,
			}

			fmt.Println(message)

			messagesjon,_ := json.Marshal(message)

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
			fmt.Printf("Response time: %s, data: %s \n", time.Now().Format("2006-01-02 15:04:05"),recvStr)
		}
	}()
}