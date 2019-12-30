package tcp

import (
	"encoding/json"
	"fmt"
	"free-im/dao"
	"free-im/model"
	"net"
)

// client auth handle
func ClientAuth(ctx *Context, mp *model.MessagePackage) {
	//认证 ctx.Message.AccessToken	&& ctx.Message.UserID
	auth := model.AuthMessage{}
	json.Unmarshal(mp.Content, &auth)
	if auth.UserID == "" || auth.AccessToken == "" || auth.DeviceID == "" {
		return
	}
	ctx.IsAuth = true
	ctx.UserID = auth.UserID
	ctx.DeviceID = auth.DeviceID

	//加入连接集合
	if model.SocketConnPool[ctx.UserID] == nil {
		model.SocketConnPool[ctx.UserID] = make(map[string] net.Conn)
	}
	model.SocketConnPool[ctx.UserID][ctx.DeviceID] = ctx.TcpConn

	//ctx.Response(model.Response{
	//	Code: 0,
	//	Msg: "认证成功",
	//})
}

//client send message handle
func ClientSendMessage(ctx *Context, mp *model.MessagePackage) {
	fmt.Println(string(mp.Content))

	redisconn := dao.NewRedis()
	defer redisconn.Close()
	//判断是否认证 auth
	if ctx.IsAuth == false {
		//ctx.ConnSocket.Close()
		return
	}
	//字段验证 code ... //

	// <- chan
	//ctx.InChan <- &ctx.Message

	// 获取消息ID
	message_id,_ := redisconn.Do("HINCRBY", "hash_im_chatroom_message_id", mp.ChatroomId, 1)
	// 存储消息
	store_message,_ := json.Marshal(mp)
	redisconn.Do("ZADD", "sorted_set_im_chatroom_message_record_"+mp.ChatroomId, message_id.(int64), store_message)
	// 查询聊天室成员
	members,_ := redisconn.Do("SMEMBERS", "set_im_chatroom_member_"+mp.ChatroomId)
	// 给聊天室全员发送消息
	send_message,_ := json.Marshal(mp)
	for _, v := range members.([]interface {}) {
		UserID := string(v.([]uint8))
		if UserID == ctx.UserID {		//消息回执
			for  k, vo := range model.SocketConnPool[ctx.UserID] {
				if k == ctx.DeviceID {
					continue
				}
				vo.Write(send_message)
			}
			continue
		}
		//给聊天室成员发送消息
		for  _, vo := range model.SocketConnPool[UserID] {
			vo.Write(send_message)
		}
	}
}

//client pull message handle
func ClientPullMessage(ctx *Context) {

}


