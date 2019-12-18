package service

import (
	"encoding/json"
	"fmt"
	"free-im/model"
	"net"
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

	//加入连接集合
	if model.SocketConnPool[ctx.Auth.UserID] == nil {
		model.SocketConnPool[ctx.Auth.UserID] = make(map[string] net.Conn)
	}
	model.SocketConnPool[ctx.Auth.UserID][ctx.Auth.DeviceID] = ctx.ConnSocket

	ctx.Response(model.Response{
		Code: 0,
		Msg: "认证成功",
	})
}

//client send message handle
func ClientSendMessage(ctx *model.Context) {
	fmt.Println(ctx.Message)
	//判断是否认证 auth
	if ctx.Auth.IsAuth == false {
		//ctx.ConnSocket.Close()
		return
	}
	//字段验证 code ... //

	// <- chan
	//ctx.InChan <- &ctx.Message

	//消息处理
	res,_ := ctx.RedisConn.Do("SMEMBERS", "set_im_chatroom_member_"+ctx.Message.ChatroomId)
	for _, v := range res.([]interface {}) {
		if string(v.([]uint8)) == ctx.Auth.UserID {		//消息回执
			for  _, vo := range model.SocketConnPool[ctx.Auth.UserID] {
				res,_ := json.Marshal(ctx.Message)
				vo.Write(res)
			}
			continue
		}
		//给聊天室全员发送消息
		for  _, vo := range model.SocketConnPool[string(v.([]uint8))] {
			res,_ := json.Marshal(ctx.Message)
			vo.Write(res)
		}
	}
}

//client pull message handle
func ClientPullMessage(ctx *model.Context) {

}


