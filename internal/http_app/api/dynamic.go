package api

import (
	"encoding/json"
	"free-im/internal/http_app/model"
	"free-im/internal/http_app/service"
	"free-im/pkg/util"
	"net/http"
	"strconv"
)

var DynamicService = new(service.DynamicService)

// 动态发布
func DynamicPublish(writer http.ResponseWriter, request *http.Request) {
	m := model.Dynamic{}
	json.NewDecoder(request.Body).Decode(&m)
	if _, err := DynamicService.Create(m); err != nil {
		util.RespFail(writer, "系统忙, 稍后再试")
		return
	}
	util.RespOk(writer, nil, "发布成功")
}

// 动态列表
func DynamicList(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]string)
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	info := make(map[string]interface{})
	page, _ := strconv.Atoi(formData["page"])
	perpage, _ := strconv.Atoi(formData["perpage"])

	list, total := DynamicService.DynamicList(page, perpage)
	info["list"] = list
	info["total"] = total
	util.RespOk(writer, info, "ok")
	return
}
