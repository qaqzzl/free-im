package v1

import (
	"encoding/json"
	"fmt"
	"free-im/server/library/cache/redis"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

// 通过会员ID 获取 聊天室ID
func MemberIdGetChatroomId(writer http.ResponseWriter, request *http.Request) {
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
		rconn.Do("HSET", "hash_im_chatroom_member_id_get_chatroom_id", field, chatroom_id)						//创建聊天室
	} else {
		chatroom_id = string( res.([]uint8) )
	}

	requestBody := fmt.Sprintf(`{
"chatroom_id":"%s",
"status": "%s",
"code": %d
}`,chatroom_id, "ok",0)
	rconn.Close()
	writer.Write([]byte(requestBody))
}

// HTTP 获取好友列表

// HTTP 获取聊天室列表

// HTTP 获取群聊列表

// HTTP 添加好友

// HTTP 删除好友