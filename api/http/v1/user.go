package v1

import (
	"encoding/json"
	"free-im/service"
	"free-im/util"
	"net/http"
)

var UserService service.UserService

// 获取个人信息
func GetUserInfo(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	info := UserService.GetMemberInfo(formData["uid"].(string))
	util.RespOk(writer, info, "")
}

// 添加好友
func AddFriend(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	var (
		err error
		info map[string]string
	)
	if info,err = UserService.AddFriend(formData["uid"].(string), formData["friend_id"].(string)); err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	}
	util.RespOk(writer, info, "")
}


// 删除好友
func DelFriend(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	if err := UserService.DelFriend(formData["uid"].(string), formData["friend_id"].(string)); err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	}
	util.RespOk(writer, nil, "删除成功")
}

// 好友申请列表
func FriendApplyList(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	apply_list,err := UserService.FriendApplyList(formData["uid"].(string))
	if err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	}
	ret := make(map[string]interface{})
	ret["apply_list"] = apply_list
	util.RespOk(writer, ret, "")
}