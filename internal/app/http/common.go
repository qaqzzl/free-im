package http

import (
	"encoding/json"
	"fmt"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/util"
	"free-im/pkg/util/id"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"net/http"
	"strconv"
	"time"
)

// 获取七牛上传token
// https://developer.qiniu.com/kodo/manual/1206/put-policy
func GetQiniuUploadToken(writer http.ResponseWriter, request *http.Request) {
	saveKeyPrefix := "test"

	accessKey := "qW7rPngWLk8Nl3MQfehQ_G5ELAZaH47Dej50Dj7k"
	secretKey := "cN5unz025wgnfHJ_Ck3iBjpLUoByXnUVB8Uu4P1g"
	putPolicy := storage.PutPolicy{
		Scope: "free-im",
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
	ret["domain"] = "http://free-im-qn.qaqzz.com/"
	ret["message"] = "获取成功"
	ret["code"] = "0"
	util.RespOk(writer, ret, "")
}

func GetMessageId(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	ret := make(map[string]interface{})
	chatroom_id, _ := strconv.Atoi(formData["chatroom_id"].(string))
	ret["message_id"] = id.MessageID.GetId(int64(chatroom_id), pbs.ChatroomType_Single)
	logger.Logger.Info(ret["message_id"].(string))
	util.RespOk(writer, ret, "")
}
