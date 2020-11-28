package http

import (
	"encoding/json"
	"free-im/pkg/util"
	"net/http"
)

// app 最新版本获取
func AppNewVersionGet(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	info := make(map[string]interface{})
	switch formData["client_type"] {
	case "android":
		info["version_code"] = 1
		info["version_name"] = "1.0.0"
		info["version_download"] = "https://cdn.qaqzz.com/free-im-v1.0.1.apk"
		info["version_description"] = "修复已知bug\n修复部分机型无法接受消息"
	case "ios":
		info["version_code"] = 1
		info["version_name"] = "1.0.0"
		info["version_download"] = "http://freeim.qaqzz.com"
		info["version_description"] = "修复已知bug\n修复部分机型无法接受消息"
	}
	util.RespOk(writer, info, "ok")
}
