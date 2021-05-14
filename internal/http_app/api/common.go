package api

import (
	"encoding/json"
	"fmt"
	"free-im/pkg/http"
	"free-im/pkg/id"
	"free-im/pkg/protos/pbs"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"strconv"
	"time"
)

// 获取七牛上传token
// https://developer.qiniu.com/kodo/manual/1206/put-policy
func GetQiniuUploadToken(c *gin.Context) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(c.Request.Body).Decode(&formData)
	var (
		scope  string
		domain string
	)
	if formData["type"] == "private" {
		scope = "free-im-private"
		domain = "http://free-im-private-qn.qaqzz.com/"
	} else if formData["type"] == "public" {
		scope = "free-im"
		domain = "http://free-im-qn.qaqzz.com/"
	}

	saveKeyPrefix := "test"
	accessKey := "qW7rPngWLk8Nl3MQfehQ_G5ELAZaH47Dej50Dj7k"
	secretKey := "cN5unz025wgnfHJ_Ck3iBjpLUoByXnUVB8Uu4P1g"
	putPolicy := storage.PutPolicy{
		Scope: scope,
		//CallbackURL:      "http://api.example.com/qiniu/upload/callback",
		CallbackBody:     `{"key":"$(key)","hash":"$(etag)","mimeType":"$(mimeType)","imageInfo":$(imageInfo),"ext":"$(ext)"}`,
		CallbackBodyType: "application/json",
		FsizeLimit:       20971520,                         //上传大小限制20M
		ForceSaveKey:     true,                             //强制使用服务端命名
		SaveKey:          saveKeyPrefix + "/$(etag)$(ext)", //强制使用服务端命名
		DetectMime:       1,                                // 使用七牛检查 mime
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	ret := make(map[string]string)
	ret["token"] = upToken
	ret["expires"] = strconv.FormatInt(time.Now().Unix()+3600, 10)
	ret["domain"] = domain
	ret["message"] = "获取成功"
	ret["code"] = "0"
	http.RespOk(c, ret, "")
}

// 发送登录短信验证码
func SendLoginSms(c *gin.Context) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(c.Request.Body).Decode(&formData)
	if err := CommonService.SendSms(formData["phone"].(string), formData["type"].(string)); err != nil {
		http.RespFail(c, err.Error())
		return
	}
	http.RespOk(c, nil, "短信验证码发送成功")
}

// * 获取消息ID
func GetMessageId(c *gin.Context) {
	var req struct {
		ChatroomID int64 `json:"chatroom_id"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	fmt.Println(strconv.Itoa(int(req.ChatroomID)))
	ret := make(map[string]interface{})
	ret["message_id"] = id.MessageID.GetId(req.ChatroomID, pbs.ChatroomType_Single)
	http.RespOk(c, ret, "")
}
