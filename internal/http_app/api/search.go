package api

import "C"
import (
	"free-im/pkg/http"
	"github.com/gin-gonic/gin"
)

// 搜索好友 | 昵称, ID
func SearchMember(c *gin.Context) {
	var req struct {
		Search string `json:"search"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	// 搜索好友列表
	search_list, err := UserService.SearchMember(http.GetUid(c), req.Search)
	if err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	http.RespOk(c, search_list, "")
}

// 搜索群组 | 群名, ID
func SearchGroup(c *gin.Context) {
	var req struct {
		Search string `json:"search"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	// 搜索群组列表
	search_list, err := ChatRoomService.SearchGroup(req.Search)
	if err != nil {
		http.RespFail(c, "系统繁忙")
		return
	}
	//ret := make(map[string]interface{})
	//ret["search_list"] = search_list
	http.RespOk(c, search_list, "")
}
