package api

import (
	"free-im/internal/http_app/model"
	"free-im/internal/http_app/service"
	"free-im/pkg/http"
	"github.com/gin-gonic/gin"
)

var DynamicService = new(service.DynamicService)

// 动态发布
func DynamicPublish(c *gin.Context) {
	m := model.Dynamic{}
	if http.ReqBin(c, &m) != nil {
		return
	}
	if _, err := DynamicService.Create(m); err != nil {
		http.RespFail(c, "系统忙, 稍后再试")
		return
	}
	http.RespOk(c, nil, "发布成功")
}

// 动态列表
func DynamicList(c *gin.Context) {
	var req struct {
		Page    int `json:"page"`
		Perpage int `json:"perpage"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	info := make(map[string]interface{})

	total, list := DynamicService.DynamicList(req.Page, req.Perpage)
	info["list"] = list
	info["total"] = total
	http.RespOk(c, info, "ok")
	return
}
