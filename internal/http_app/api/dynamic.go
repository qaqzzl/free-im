package api

import (
	"encoding/json"
	"free-im/internal/http_app/model"
	"free-im/internal/http_app/service"
	"free-im/pkg/http"
	"github.com/gin-gonic/gin"
	"strconv"
)

var DynamicService = new(service.DynamicService)

// 动态发布
func DynamicPublish(c *gin.Context) {
	m := model.Dynamic{}
	json.NewDecoder(c.Request.Body).Decode(&m)
	if _, err := DynamicService.Create(m); err != nil {
		http.RespFail(c, "系统忙, 稍后再试")
		return
	}
	http.RespOk(c, nil, "发布成功")
}

// 动态列表
func DynamicList(c *gin.Context) {
	// 初始化请求变量结构
	formData := make(map[string]string)
	// 调用json包的解析，解析请求body
	json.NewDecoder(c.Request.Body).Decode(&formData)
	info := make(map[string]interface{})
	page, _ := strconv.Atoi(formData["page"])
	perpage, _ := strconv.Atoi(formData["perpage"])

	list, total := DynamicService.DynamicList(page, perpage)
	info["list"] = list
	info["total"] = total
	http.RespOk(c, info, "ok")
	return
}
