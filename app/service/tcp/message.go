package tcp

import (
	"encoding/json"
	"fmt"
	"free-im/app/model"
	"free-im/dao"
	"net"
	"strconv"
	"time"
)

// client auth handle
func ClientAuth(ctx *Context, message model.AuthMessage) {
	fmt.Println("进入认证方法")
	//认证
	if message.UserID == "" || message.AccessToken == "" || message.DeviceID == "" {
		return
	}
	ctx.IsAuth = true
	ctx.UserID = message.UserID
	ctx.DeviceID = message.DeviceID

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
func ClientSendMessage(ctx *Context, message model.MessagePackage) {
	fmt.Println("进入消息接受处理程序", ctx.UserID)
	message.UserId = ctx.UserID
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
	origin_message_id,_ := redisconn.Do("HINCRBY", "hash_im_chatroom_message_id", message.ChatroomId, 1)
	message_id := origin_message_id.(int64)

	// 存储消息
	store_message,_ := json.Marshal(message)
	redisconn.Do("ZADD", "sorted_set_im_chatroom_message_record:"+message.ChatroomId, message_id, store_message)

	// 查询聊天室成员
	members,_ := redisconn.Do("SMEMBERS", "set_im_chatroom_member:"+message.ChatroomId)
	fmt.Println("查询聊天室成员",members)
	// 给聊天室全员发送消息
	server_send_message := model.ServerSendMessagePackage{
		Code:message.Code,
		ChatroomId:message.ChatroomId,
		Content:message.Content,
		ClientMessageId:message.MessageId,
		ServerMessageId:strconv.FormatInt(message_id,10),
		UserId:message.UserId,
		MessageSendTime:time.Now().Unix(),
	}
	send_message,_ := json.Marshal(server_send_message)
	other_send_message := model.ActionMessageSend + string(send_message);			// 消息
	own_send_message := model.ActionMessageSendACK + string(send_message);		// 消息回执
	for _, v := range members.([]interface {}) {
		UserID := string(v.([]uint8))
		if UserID == ctx.UserID {		// 其他设备消息同步
			for  k, vo := range model.SocketConnPool[ctx.UserID] {
				if k == ctx.DeviceID {	// 消息回执
					fmt.Println("消息回执:", own_send_message)
					vo.Write([]byte(own_send_message))
					continue
				}
				// 同步其他客户端
				vo.Write([]byte(other_send_message))
			}
			continue
		}
		//给聊天室成员发送消息
		fmt.Println("消息:", other_send_message)
		for  _, vo := range model.SocketConnPool[UserID] {
			vo.Write([]byte(other_send_message))
		}
	}
}

//client pull message handle
func ClientPullMessage(ctx *Context) {

}


