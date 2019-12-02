package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"free-im/client/model"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:1208")
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
		return
	}

	TestHandler(conn)


	<-time.Tick(time.Second * 10000)
}


//测试脚本
func TestHandler(c net.Conn) {

	for i:=1; i<10; i++ {
		<- time.Tick(time.Millisecond * 10)
	go func() {		//发送消息
		//登录
		fmt.Println(i)
		//获取聊天室ID
		requestBody := fmt.Sprintf(`{
"user_id":"%s",
"access_token": "%s",
"member_id": "%s"
}`,
strconv.Itoa(i), "access_token","1")
		//fmt.Println(requestBody)

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
		chatroom_id := formData["chatroom_id"].(string)		//聊天室ID
		fmt.Println(chatroom_id)

		//socket连接认证
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
		//循环发送消息
		for {
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


	}

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

//获取聊天室ID
func getChatroomId(c net.Conn, member_id string) {
	content := make(map[string]string)
	content["member_id"] = member_id
	message := model.MessagePackage{
		Motion:	12,
		Content:content,
	}

	messagesjon,_ := json.Marshal(message)
	fmt.Println(messagesjon)

	c.Write(messagesjon)
}