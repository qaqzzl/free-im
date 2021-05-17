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
	var req struct {
		FriendID int64 `json:"friend_id"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	var (
		chatroom_id string
		err         error
	)
	if chatroom_id, err = ChatRoomService.FriendIdGetChatroomId(http.GetUid(c), req.FriendID); err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	ret := make(map[string]string)
	ret["chatroom_id"] = chatroom_id
	http.RespOk(c, ret, "")
}

// 通过聊天室ID获取聊天室基础信息（头像，名称）
func GetChatroomAvatarNameByChatRoomID(c *gin.Context) {
	var req struct {
		ChatroomID int64 `json:"chatroom_id"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	if res, err := ChatRoomService.GetChatroomBaseInfo(req.ChatroomID, http.GetUid(c)); err != nil {
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
	var req struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	var (
		group_id int64
		err      error
	)
	if group_id, err = ChatRoomService.CreateGroup(http.GetUid(c), model.Group{
		Name:   req.Name,
		Avatar: req.Avatar,
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
	var req struct {
		GroupID int64  `json:"group_id"`
		Remark  string `json:"remark"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}

	if ret, err := ChatRoomService.JoinGroup(http.GetUid(c), req.GroupID, req.Remark); err != nil {
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
	group_list, err := ChatRoomService.MemberGroupList(http.GetUid(c))
	if err != nil {
		http.RespFail(c, err.Error())
		return
	}
	ret := make(map[string]interface{})
	ret["group_list"] = group_list
	http.RespOk(c, ret, "")
}
