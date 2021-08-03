package api

import (
	"free-im/pkg/http"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// http_app 最新版本获取
func AppNewVersionGet(c *gin.Context) {
	info := make(map[string]interface{})
	switch http.GetClientType(c) {
	case "ios":
		info["version_code"] = 1
		info["is_must"] = 1
		info["version_name"] = "1.0.0"
		info["version_download"] = ""
		info["version_download_page"] = "https://www.pgyer.com/freeim"
		info["version_description"] = ""
	default:
		info["version_code"] = viper.GetString("AndroidVersionCode")
		info["is_must"] = 1 // 是否必须更新
		info["version_name"] = viper.GetString("AndroidVersionName")
		info["version_download"] = "http://cdn.qaqzz.com/free-im-1.0.2.apk"
		info["version_download_page"] = viper.GetString("AndroidVersionPage")
		info["version_description"] = viper.GetString("AndroidVersionDesc")
	}
	http.RespOk(c, info, "ok")
}
