package v1

import (
	"encoding/json"
	"free-im/model"
	"free-im/service"
	"free-im/util"
	"net/http"
)

var ChatRoomService service.ChatRoomService

// 通过好友ID 获取 聊天室ID
func FriendIdGetChatroomId(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	var (
		chatroom_id string
		err error
	)
	if chatroom_id, err = ChatRoomService.FriendIdGetChatroomId(formData["uid"].(string), formData["friend_id"].(string)); err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	}
	ret := make(map[string]string)
	ret["chatroom_id"] = chatroom_id
	util.RespOk(writer, ret, "")
}

// 聊天室列表 -- 未完成
func ChatroomList(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	ChatRoomService.ChatroomList(formData["uid"].(string))
}

// 创建群组
func CreateGroup(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	ChatRoomService.CreateGroup(formData["uid"].(string), model.Group{
		Name: formData["name"].(string),
		Avatar: formData["avatar"].(string),
	})

}

// 加入群组
func AddGroup(writer http.ResponseWriter, request *http.Request) {

}

// 退出群组
func OutGroup(writer http.ResponseWriter, request *http.Request) {

}

// 我的群组列表
func MyGroupList(writer http.ResponseWriter, request *http.Request) {

}
