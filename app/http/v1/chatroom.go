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

	var (
		group_id string
		err error
	)
	if group_id, err = ChatRoomService.CreateGroup(formData["uid"].(string), model.Group{
		Name: formData["name"].(string),
		Avatar: formData["avatar"].(string),
	}); err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	}

	ret := make(map[string]string)
	ret["group_id"] = group_id
	util.RespOk(writer, ret, "")
}

// 加入群组
func AddGroup(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)

	if ret,err := ChatRoomService.AddGroup(formData["uid"].(string), formData["group_id"].(string), formData["remark"].(string)); err != nil {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, ret, "")
	}

}

// 退出群组
func OutGroup(writer http.ResponseWriter, request *http.Request) {

}

// 我的群组列表
func MyGroupList(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)

	group_list, err := ChatRoomService.MyGroupList(formData["uid"].(string))
	if err != nil {
		util.RespFail(writer, err.Error())
		return
	}
	ret := make(map[string]interface{})
	ret["group_list"] = group_list
	util.RespOk(writer, ret, "")
}
