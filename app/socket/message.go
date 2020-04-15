package socket

import (
	"encoding/json"
	"fmt"
	"free-im/app/model"
	"free-im/dao"
	"github.com/orcaman/concurrent-map"
	"net"
	"strconv"
	"time"
)

type Socket interface {
	ClientAuth(message model.AuthMessage)
	ClientMessage(message model.MessagePackage)
	Close()
}

type Message struct {
	ctx *Context
}

// client auth handle
func (message Message) ClientAuth(m model.AuthMessage) {
	fmt.Println("进入认证方法",m)
	//认证
	if m.UserID == "" || m.AccessToken == "" || m.DeviceID == "" || m.ClientType == "" || m.DeviceType == "" {
		return
	}
	message.ctx.IsAuth = true
	message.ctx.UserID = m.UserID
	message.ctx.DeviceID = m.DeviceID
	message.ctx.ClientType = m.ClientType
	message.ctx.DeviceType = m.DeviceType
	clientDevice := model.ClientDevice{
		DeviceID:m.DeviceID,
		ClientType:m.ClientType,
		Conn:message.ctx.TcpConn,
	}
	//加入连接集合
	if tmp, ok := model.SocketConnPool.Get(message.ctx.UserID); ok {
		device_map := tmp.(cmap.ConcurrentMap)
		// 判断连接是否存在相同设备
		for k,v := range device_map.Items() {
			if k == m.DeviceType {	// 如果有同类型的设备登录了 ,通知其设备下线
				device := v.(model.ClientDevice)
				if device.DeviceID != m.DeviceID {
					// 通知其设备下线 code...
					fmt.Println("通知其设备下线", device.DeviceID)
					// 关闭连接
				}
			}
		}
		device_map.Set(m.DeviceType, clientDevice)
		model.SocketConnPool.Set(message.ctx.UserID, device_map)
	} else  {
		device_map := cmap.New()
		device_map.Set(m.DeviceType, clientDevice)
		model.SocketConnPool.Set(message.ctx.UserID, device_map)
	}

	//ctx.Response(model.Response{
	//	Code: 0,
	//	Msg: "认证成功",
	//})
}

//client send message handle
func (message Message) ClientMessage(m model.MessagePackage) {
	fmt.Println("进入消息接受处理程序", message.ctx.UserID)
	m.UserId = message.ctx.UserID
	redisconn := dao.NewRedis()
	defer redisconn.Close()
	//判断是否认证 auth
	if message.ctx.IsAuth == false {
		//ctx.ConnSocket.Close()
		return
	}
	//字段验证 code ... //

	// <- chan
	//ctx.InChan <- &ctx.Message

	// 获取消息ID
	origin_message_id,_ := redisconn.Do("HINCRBY", "hash_im_chatroom_message_id", m.ChatroomId, 1)
	message_id := origin_message_id.(int64)

	// 存储消息
	store_message,_ := json.Marshal(message)
	redisconn.Do("ZADD", "sorted_set_im_chatroom_message_record:"+m.ChatroomId, message_id, store_message)

	// 查询聊天室成员
	members,_ := redisconn.Do("SMEMBERS", "set_im_chatroom_member:"+m.ChatroomId)
	// 给聊天室全员发送消息
	server_send_message := model.ServerSendMessagePackage{
		Code:m.Code,
		ChatroomId:m.ChatroomId,
		Content:m.Content,
		ClientMessageId:m.MessageId,
		ServerMessageId:strconv.FormatInt(message_id,10),
		UserId:m.UserId,
		MessageSendTime:time.Now().Unix(),
	}
	send_message,_ := json.Marshal(server_send_message)
	other_send_message := model.ActionMessageSend + string(send_message)			// 消息投递
	own_send_message := model.ActionMessageACK + string(send_message)		// 消息回执
	for _, v := range members.([]interface {}) {
		UserID := string(v.([]uint8))
		fmt.Println("给其他UserID发送消息",UserID)
		if UserID == message.ctx.UserID {		// 其他设备消息同步
			if tmp, ok := model.SocketConnPool.Get(message.ctx.UserID); ok {
				for  k, vo := range tmp.(cmap.ConcurrentMap).Items() {
					device := vo.(model.ClientDevice)
					if k == message.ctx.DeviceType {	// 消息回执
						fmt.Println("消息回执:", own_send_message)
						message.ctx.OutChan <- sendMessage{
							Conn:device.Conn,
							ReceiveUserID:UserID,
							Message:[]byte(own_send_message),
						}
						continue
					}
					// 同步其他客户端
					message.ctx.OutChan <- sendMessage{
						Conn:device.Conn,
						ReceiveUserID:UserID,
						Message:[]byte(other_send_message),
					}
				}
			}
			continue
		}

		//给聊天室成员发送消息
		fmt.Println("消息:", other_send_message)
		tmp, ok := model.SocketConnPool.Get(UserID)
		if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
			for  _, vo := range tmp.(cmap.ConcurrentMap).Items() {
				device := vo.(model.ClientDevice)
				message.ctx.OutChan <- sendMessage{
					Conn:device.Conn,
					ReceiveUserID:UserID,
					Message:[]byte(other_send_message),
				}
			}
		} else {
			fmt.Println("设备未在线 , 未读消息写入redis")
			redisconn.Do("LPUSH", "list_message_send_failure:"+UserID, other_send_message)
		}

	}
}

func (message Message) Close() {
	if user_map, ok := model.SocketConnPool.Get(message.ctx.UserID); ok {
		fmt.Println(message.ctx.DeviceType)
		user_map.(cmap.ConcurrentMap).Remove(message.ctx.DeviceType)
		model.SocketConnPool.Set(message.ctx.UserID,user_map)
	}
	message.ctx.IsConnStatus = false
	message.ctx.TcpConn.Close()
}

//client pull message handle
func (message Message) SendResponse(conn net.Conn, ReceiveUserID string,m []byte) (n int, err error) {
	redisconn := dao.NewRedis()
	defer redisconn.Close()
	if n, err = message.ctx.Write(conn,m); err != nil {
		fmt.Println("消息发送失败 , 写入redis")
		// 消息发送失败 , 写入redis
		redisconn.Do("LPUSH", "list_message_send_failure:"+ReceiveUserID, m)
	}
	return n,err
}

func TcpClose(ctx *Context) {
	ctx.IsConnStatus = false
	ctx.TcpConn.Close()
}




