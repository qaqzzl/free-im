package api

import (
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
		chatroom_id int64
		err         error
	)
	if chatroom_id, err = ChatRoomService.FriendIdGetChatroomId(http.GetUid(c), req.FriendID); err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	ret := make(map[string]interface{})
	ret["chatroom_id"] = chatroom_id
	http.RespOk(c, ret, "")
}

// 获取聊天室信息
func GetChatroomBaseInfo(c *gin.Context) {
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
	ChatRoomService.ChatroomList(http.GetUid(c))
}

// 创建群组
func CreateGroup(c *gin.Context) {
	var req struct {
		Name       string  `json:"name"`
		Avatar     string  `json:"avatar"`
		Desc       string  `json:"desc"`
		MemberList []int64 `json:"member_list"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	group := model.Group{
		Name:   req.Name,
		Avatar: req.Avatar,
		Desc:   req.Desc,
	}
	if err := ChatRoomService.CreateGroup(http.GetUid(c), &group, req.MemberList); err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	http.RespOk(c, group, "")
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
	http.RespOk(c, group_list, "")
}

// 群组信息
func GroupInfo(c *gin.Context) {
	var req struct {
		GroupID int64 `json:"group_id"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	info, err := ChatRoomService.GroupInfo(req.GroupID)
	if err != nil {
		http.RespFail(c, err.Error())
		return
	}
	http.RespOk(c, info, "")
}

// 群组成员
func GroupMember(c *gin.Context) {
	var req struct {
		GroupID int64 `uri:"group_id" binding:"required"`
	}
	if err := http.ShouldBindUri(c, &req); err != nil {
		return
	}

	list, err := ChatRoomService.GroupMember(http.GetUid(c), req.GroupID)
	if err != nil {
		http.RespFail(c, err.Error())
		return
	}
	http.RespOk(c, list, "")
}

// 添加群成员
func AddGroupMember(c *gin.Context) {
	var req struct {
		GroupID    int64   `json:"group_id"`
		MemberList []int64 `json:"member_list"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	// todo 验证数据 code。。。
	if group_member, err := ChatRoomService.AddGroupMember(http.GetUid(c), req.GroupID, req.MemberList); err != nil {
		http.RespFail(c, "系统繁忙")
		return
	} else {
		http.RespOk(c, group_member, "")
	}
}
