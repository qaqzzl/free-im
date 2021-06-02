package api

import (
	"free-im/pkg/http"
	"github.com/gin-gonic/gin"
)

// http_app 最新版本获取
func AppNewVersionGet(c *gin.Context) {
	info := make(map[string]interface{})
	switch http.GetClientType(c) {
	case "android":
		info["version_code"] = 3
		info["is_must"] = 1 // 是否必须更新
		info["version_name"] = "2.0.0"
		info["version_download"] = "https://cdn.qaqzz.com/free-im-v2.0.0.apk"
		info["version_download_page"] = "https://www.pgyer.com/freeim"
		info["version_description"] = "全新版本\n新增群聊功能\n修复已知BUG"
	case "ios":
		info["version_code"] = 1
		info["is_must"] = 1
		info["version_name"] = "1.0.0"
		info["version_download"] = ""
		info["version_download_page"] = "https://www.pgyer.com/freeim"
		info["version_description"] = ""
	default:
		info["version_code"] = 3
		info["is_must"] = 1 // 是否必须更新
		info["version_name"] = "2.0.0"
		info["version_download"] = ""
		info["version_download_page"] = "https://www.pgyer.com/freeim"
		info["version_description"] = "全新版本\n新增群聊功能\n修复已知BUG"
	}
	http.RespOk(c, info, "ok")
}
