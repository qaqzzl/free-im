package api

import (
	"encoding/json"
	"free-im/pkg/http"
	"github.com/gin-gonic/gin"
)

// http_app 最新版本获取
func AppNewVersionGet(c *gin.Context) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(c.Request.Body).Decode(&formData)
	info := make(map[string]interface{})
	switch formData["client_type"] {
	case "android":
		info["version_code"] = 2
		info["version_name"] = "1.1.0"
		info["version_download"] = "https://cdn.qaqzz.com/free-im-v1.1.0.apk"
		info["version_description"] = "修复已知bug"
	case "ios":
		info["version_code"] = 2
		info["version_name"] = "1.1.0"
		info["version_download"] = "http://freeim.qaqzz.com"
		info["version_description"] = "修复已知bug"
	}
	http.RespOk(c, info, "ok")
}
