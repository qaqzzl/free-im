package api

import (
	"free-im/internal/http_app/model"
	"free-im/internal/http_app/service"
	"free-im/pkg/http"
	"github.com/gin-gonic/gin"
)

var UserService service.UserService

// 获取个人信息
func GetMemberInfo(c *gin.Context) {
	info, _ := UserService.GetMemberInfo(http.GetUid(c))
	http.RespOk(c, info, "")
}

// 修改个人信息
func UpdateMemberInfo(c *gin.Context) {
	var req struct {
		Nickname  string `json:"nickname"`
		Gender    string `json:"Gender"`
		Birthdate int    `json:"Birthdate"`
		Avatar    string `json:"avatar"`
		Signature string `json:"signature"`
		City      string `json:"city"`
		Province  string `json:"province"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	uid := http.GetUid(c)
	var m = model.UserMember{
		Nickname:  req.Nickname,
		Gender:    req.Gender,
		Birthdate: req.Birthdate,
		Avatar:    req.Avatar,
		Signature: req.Signature,
		City:      req.City,
		Province:  req.Province,
	}
	if err := UserService.UpdateMemberInfo(uid, m); err != nil {
		http.RespFail(c, err.Error())
		return
	}
	http.RespOk(c, nil, "修改成功")
}

// 添加好友
func AddFriend(c *gin.Context) {
	var req struct {
		FriendID int64  `json:"friend_id"`
		Remark   string `json:"remark"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	var (
		err  error
		info map[string]string
	)
	uid := http.GetUid(c)
	if req.FriendID == uid {
		http.RespFail(c, "不可以添加自己为好友哦")
		return
	}
	if info, err = UserService.AddFriend(uid, req.FriendID, req.Remark); err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	http.RespOk(c, info, "")
}

// 删除好友
func DelFriend(c *gin.Context) {
	var req struct {
		FriendID int64 `json:"friend_id"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	if err := UserService.DelFriend(http.GetUid(c), req.FriendID); err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	http.RespOk(c, nil, "删除成功")
}

// 好友申请列表
func FriendApplyList(c *gin.Context) {
	apply_list, err := UserService.FriendApplyList(http.GetUid(c))
	if err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	ret := make(map[string]interface{})
	ret["apply_list"] = apply_list
	http.RespOk(c, ret, "")
}

// 好友申请同意/拒绝操作
func FriendApplyAction(c *gin.Context) {
	var req struct {
		ID     int64 `json:"id"`
		Action int   `json:"action"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	ret, err := UserService.FriendApplyAction(req.ID, http.GetUid(c), req.Action)
	if err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	http.RespOk(c, ret, "")
}

// 好友列表
func FriendList(c *gin.Context) {
	apply_list, err := UserService.FriendList(http.GetUid(c))
	if err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	http.RespOk(c, apply_list, "")
}

// 他人主页(用户基本信息)
func OthersHomeInfo(c *gin.Context) {
	var req struct {
		MemberId int64 `json:"member_id"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	info, _ := UserService.OthersHomeInfo(http.GetUid(c), req.MemberId)
	http.RespOk(c, info, "")
}

/*
2021-06-01 23:24:43.551 ERROR   api/user.go:139 他人主页(用户基本信息)map[avatar:http://free-im-qn.qaqzz.com/test/FqrLCxJd1GkG2ObKyS09NurfBVRX.jpg birthdate:852048om_id:0 city: gender:m is_friend:no member_id:1 nickname:龙猫 province: signature:需要个android帮我]

*/
