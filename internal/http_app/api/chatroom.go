package api

import (
	"encoding/json"
	"free-im/internal/http_app/model"
	"free-im/internal/http_app/service"
	"free-im/pkg/http"
	"github.com/gin-gonic/gin"
)

var ChatRoomService service.ChatRoomService

// 通过好友ID 获取 聊天室ID
func FriendIdGetChatroomId(c *gin.Context) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(c.Request.Body).Decode(&formData)
	var (
		chatroom_id string
		err         error
	)

	if chatroom_id, err = ChatRoomService.FriendIdGetChatroomId(formData["uid"].(string), formData["friend_id"].(string)); err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	ret := make(map[string]string)
	ret["chatroom_id"] = chatroom_id
	http.RespOk(c, ret, "")
}

// 通过聊天室ID获取聊天室基础信息（头像，名称）
func GetChatroomAvatarNameByChatRoomID(c *gin.Context) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(c.Request.Body).Decode(&formData)
	user_id := formData["uid"].(int)
	chatroom_id := formData["chatroom_id"].(string)
	chatroom_type := "1"
	if res, err := ChatRoomService.GetChatroomBaseInfo(chatroom_id, chatroom_type, int64(user_id)); err != nil {
		http.RespFail(c, err.Error())
	} else {
		http.RespOk(c, res, "")
	}
}

// 聊天室列表
func ChatroomList(c *gin.Context) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(c.Request.Body).Decode(&formData)
	ChatRoomService.ChatroomList(formData["uid"].(string))
}

// 创建群组
func CreateGroup(c *gin.Context) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(c.Request.Body).Decode(&formData)

	var (
		group_id int64
		err      error
	)
	if group_id, err = ChatRoomService.CreateGroup(int64(formData["uid"].(int)), model.Group{
		Name:   formData["name"].(string),
		Avatar: formData["avatar"].(string),
	}); err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}

	ret := make(map[string]interface{})
	ret["group_id"] = group_id
	http.RespOk(c, ret, "")
}

// 加入群组
func AddGroup(c *gin.Context) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(c.Request.Body).Decode(&formData)

	if ret, err := ChatRoomService.JoinGroup(int64(formData["uid"].(int)), formData["group_id"].(string), formData["remark"].(string)); err != nil {
		http.RespFail(c, err.Error())
	} else {
		http.RespOk(c, ret, "")
	}

}

// 退出群组
func OutGroup(c *gin.Context) {

}

// 我的群组列表
func MyGroupList(c *gin.Context) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(c.Request.Body).Decode(&formData)

	group_list, err := ChatRoomService.MemberGroupList(formData["uid"].(uint))
	if err != nil {
		http.RespFail(c, err.Error())
		return
	}
	ret := make(map[string]interface{})
	ret["group_list"] = group_list
	http.RespOk(c, ret, "")
}
