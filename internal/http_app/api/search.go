package api

import "C"
import (
	"encoding/json"
	"free-im/pkg/http"
	"github.com/gin-gonic/gin"
)

// 搜索好友 | 昵称, ID
func SearchMember(c *gin.Context) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(c.Request.Body).Decode(&formData)
	// 搜索好友列表
	search_list, err := UserService.SearchMember(formData["search"].(string))
	if err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	ret := make(map[string]interface{})
	ret["search_list"] = search_list
	http.RespOk(c, ret, "")
}

// 搜索群组 | 群名, ID
func SearchGroup(c *gin.Context) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(c.Request.Body).Decode(&formData)
	// 搜索好友列表
	search_list, err := ChatRoomService.SearchGroup(formData["search"].(string))
	if err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	ret := make(map[string]interface{})
	ret["search_list"] = search_list
	http.RespOk(c, ret, "")
}
