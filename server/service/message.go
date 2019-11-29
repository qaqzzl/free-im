package service

import (
	"free-im/server/model"
)

// client auth handle
func ClientAuth(ctx *model.Context) {
	//认证 ctx.Message.AccessToken	&& ctx.Message.UserID
	content := ctx.Message.Content.(map[string]interface{})
	if content["user_id"] == nil || content["access_token"] == nil || content["device_id"] == nil {
		return
	}
	ctx.Auth.IsAuth = true
	ctx.Auth.UserID = content["user_id"].(string)
	ctx.Auth.AccessToken = content["access_token"].(string)
	ctx.Auth.DeviceID = content["device_id"].(string)

	ctx.Response(model.Response{
		Code: 0,
		Msg: "认证成功",
	})
}

//client send message handle
func ClientSendMessage(ctx *model.Context) {
	//判断是否认证 auth
	if ctx.Auth.IsAuth == false {
		ctx.ConnSocket.Close()
		return
	}
	//字段验证 code ... //


}

//client pull message handle
func ClientPullMessage(ctx *model.Context) {

}


