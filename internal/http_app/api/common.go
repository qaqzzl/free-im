package api

import (
	"context"
	"fmt"
	"free-im/config"
	"free-im/pkg/http"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/rpc_client"
	"free-im/pkg/service/id"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"strconv"
	"time"
)

// 获取七牛上传token
// https://developer.qiniu.com/kodo/manual/1206/put-policy
func GetQiniuUploadToken(c *gin.Context) {
	var req struct {
		Type string `json:"type"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	var (
		scope  string
		domain string
	)
	switch req.Type {
	case "private":
		scope = "free-im-private"
		domain = "http://free-im-private-qn.qaqzz.com/"
	case "public":
		scope = "free-im"
		domain = "http://free-im-qn.qaqzz.com/"
	default:
		scope = "free-im"
		domain = "http://free-im-qn.qaqzz.com/"
	}
	saveKeyPrefix := "dev"
	accessKey := config.CommonConf.QiniuAccessKey
	secretKey := config.CommonConf.QiniuSecretKey
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
	var req struct {
		Phone string `json:"phone"`
		Type  string `json:"type"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	if err := CommonService.SendSms(req.Phone, req.Type); err != nil {
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

// * 发送消息
func PushMessage(c *gin.Context) {
	var req struct {
		ChatroomID   int64  `json:"chatroom_id"`
		ChatroomType int    `json:"chatroom_type"`
		Code         int    `json:"code"`
		Content      string `json:"content"`
		ID           string `json:"id"`
	}
	if http.ReqBin(c, &req) != nil {
		return
	}
	// todo 权限验证 code 。。。
	if req.ChatroomID == 0 {
		logger.Sugar.Error("聊天室ID不能为空")
		http.RespFail(c, "聊天室ID不能为空")
		return
	}

	var message_id = id.MessageID.GetId(req.ChatroomID, pbs.ChatroomType_Single)
	var m = &pbs.MsgItem{
		Code:            pbs.MessageCode(req.Code),
		ChatroomType:    pbs.ChatroomType(req.ChatroomType),
		ChatroomId:      req.ChatroomID,
		Content:         req.Content,
		MessageId:       message_id,
		UserId:          http.GetUid(c),
		DeviceID:        http.GetDeviceId(c),
		ClientType:      http.GetClientType(c),
		MessageSendTime: time.Now().Unix(),
	}
	_, _ = rpc_client.LogicInit.MessageReceive(context.TODO(), &pbs.MessageReceiveReq{
		Message: m,
	})

	ret := make(map[string]interface{})
	ret["id"] = req.ID
	ret["code"] = req.Code
	ret["chatroom_type"] = req.ChatroomType
	ret["chatroom_id"] = req.ChatroomID
	ret["content"] = req.Content
	ret["message_id"] = message_id
	ret["user_id"] = http.GetUid(c)
	ret["device_id"] = http.GetDeviceId(c)
	ret["client_type"] = http.GetClientType(c)
	ret["message_send_time"] = m.MessageSendTime
	fmt.Println(req)
	fmt.Println(ret)
	http.RespOk(c, ret, "")
}
