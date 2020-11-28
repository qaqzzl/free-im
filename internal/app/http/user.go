package http

import (
	"encoding/json"
	"free-im/internal/app/service"
	"free-im/pkg/util"
	"log"
	"net/http"
)

var UserService service.UserService

// 获取个人信息
func GetMemberInfo(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	info, _ := UserService.GetMemberInfo(formData["uid"].(string))
	util.RespOk(writer, info, "")
}

// 修改个人信息
func UpdateMemberInfo(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]string)
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	info := make(map[string]string)
	member_id := formData["uid"]
	// 判断用户昵称是否存在
	if _, is_nickname := formData["nickname"]; is_nickname {
		if formData["nickname"] == "" {
			info["message"] = "昵称不能为空"
			info["code"] = "1"
			util.RespOk(writer, info, "")
			return
		}
		if UserService.IsMemberNickname(formData["uid"], formData["nickname"]) {
			info["message"] = "昵称已经被使用"
			info["code"] = "1"
			util.RespOk(writer, info, "")
			return
		}
	}
	delete(formData, "access_token")
	delete(formData, "uid")
	if err := UserService.UpdateMemberInfo(member_id, formData); err != nil {
		log.Println(err)
		util.RespFail(writer, "系统繁忙")
		return
	}
	info["message"] = "保存成功"
	info["code"] = "0"
	util.RespOk(writer, info, "")
}

// 添加好友
func AddFriend(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	var (
		err  error
		info map[string]string
	)
	if formData["friend_id"].(string) == formData["uid"].(string) {
		util.RespFail(writer, "不可以添加自己为好友哦")
		return
	}
	if info, err = UserService.AddFriend(formData["uid"].(string), formData["friend_id"].(string), formData["remark"].(string)); err != nil {
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
	apply_list, err := UserService.FriendApplyList(formData["uid"].(string))
	if err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	}
	ret := make(map[string]interface{})
	ret["apply_list"] = apply_list
	util.RespOk(writer, ret, "")
}

// 好友申请同意/拒绝操作
func FriendApplyAction(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	ret, err := UserService.FriendApplyAction(formData["id"].(string), formData["uid"].(string), formData["action"].(string))
	if err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	}
	util.RespOk(writer, ret, "")
}

// 好友列表
func FriendList(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	apply_list, err := UserService.FriendList(formData["uid"].(string))
	if err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	}
	ret := make(map[string]interface{})
	ret["friend_list"] = apply_list

	util.RespOk(writer, ret, "")
}

// 他人主页(用户基本信息)
func OthersHomeInfo(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	info, _ := UserService.OthersHomeInfo(formData["uid"].(string), formData["member_id"].(string))
	util.RespOk(writer, info, "")
}
