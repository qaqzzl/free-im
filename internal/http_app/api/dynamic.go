package api

import (
	"encoding/json"
	"free-im/internal/http_app/dao"
	"free-im/internal/http_app/model"
	"free-im/internal/http_app/service"
	"free-im/pkg/util"
	"net/http"
	"strconv"
	"time"
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
	int_current_page, _ := strconv.Atoi(formData["current_page"])
	int_perpage, _ := strconv.Atoi(formData["perpage"])
	current_page := strconv.Itoa((int_current_page - 1) * int_perpage)
	total, _ := dao.NewMysql().Table("dynamic").Count()
	if list, err := dao.NewMysql().Table("dynamic as d").
		Join("join user_member um on um.member_id = d.member_id").
		Select("d.*,um.nickname,um.avatar,um.gender,um.birthdate").
		Order("dynamic_id desc").
		Limit(current_page + "," + formData["perpage"]).Get(); err != nil {
		util.RespFail(writer, "系统忙, 稍后再试")
		return
	} else {
		if len(list) == 0 {
			list = make([]map[string]string, 0)
		}

		for k, v := range list {
			created_at, _ := strconv.ParseInt(v["created_at"], 10, 64)
			tm := time.Unix(created_at, 0)
			list[k]["created_at"] = tm.Format("2006-01-02 15:04:05")
		}
		info["list"] = list
		info["total"] = total
		util.RespOk(writer, info, "ok")
		return
	}
}
