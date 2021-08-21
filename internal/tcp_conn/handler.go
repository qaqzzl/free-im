package tcp_conn

import (
	"context"
	"free-im/pkg/logger"
	"free-im/pkg/msg"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/rpc_client"
	"github.com/orcaman/concurrent-map"
	"strconv"
)

type handler struct{}

var Handler = new(handler)

var msgFormat = msg.MsgFormat{
	Coding: "json",
}

func (h *handler) Handler(conn *Conn, mp pbs.MsgPackage) {
	switch mp.Action {
	case pbs.Action_GetMessageID: // 获取消息ID
	case pbs.Action_Auth: // 连接认证
		h.Auth(conn, mp)
	case pbs.Action_Message: // 消息
		h.MessageReceive(conn, mp)
	case pbs.Action_MessageACK: // 消息回执
		h.MessageACK(conn, mp)
	case pbs.Action_SyncTrigger: // 消息同步
		h.SyncTrigger(conn, mp)
	case pbs.Action_Headbeat: // 心跳
		h.Headbeat(conn)
	case pbs.Action_Quit:
		conn.Close()
	default:
		logger.Sugar.Error("Unsupported command:", mp)
	}
}

// client auth handle
func (h *handler) Auth(conn *Conn, mp pbs.MsgPackage) {
	m := pbs.MsgAuth{}
	if err := msgFormat.Decode(mp.BodyData, &m); err != nil {
		logger.Sugar.Error(err)
		return
	}

	if resp, err := rpc_client.LogicInit.TokenAuth(context.TODO(), &pbs.TokenAuthReq{
		Message: &m,
	}); err != nil {
		logger.Sugar.Error(err)
		return
	} else if resp.Statu == false {
		return
	}

	conn.IsAuth = true
	conn.UserID = m.UserID
	conn.DeviceID = m.DeviceID
	conn.ClientType = m.ClientType
	conn.DeviceType = m.DeviceType
	//加入连接集合
	key := strconv.Itoa(int(conn.UserID))
	if tmp, ok := TCPServer.ServerConnPool.Get(key); ok {
		device_map := tmp.(cmap.ConcurrentMap)
		// 判断连接是否存在相同设备
		for k, v := range device_map.Items() {
			if k == m.DeviceType { // 如果有同类型的设备登录了 ,通知其设备下线
				vconn := v.(*Conn)
				if vconn.DeviceID != m.DeviceID {
					// 通知其设备下线
					msgQuit := &pbs.MsgQuit{
						Title:   "其他设备登陆通知",
						Content: "你的账号在其他设备登陆<br>如不是你本人登陆请<a href=''>修改密码</a>",
					}
					BodyData, _ := msgFormat.Encode(msgQuit)
					vconn.Write(pbs.MsgPackage{
						Version:  conn.Version,
						Action:   pbs.Action_Quit,
						BodyData: BodyData,
					})
				}
				// 关闭连接
				vconn.Close()
			}
		}
	}
	// 储存连接
	TCPServer.StoreConn(conn)
	// 发生连接认证成功消息
	conn.Write(pbs.MsgPackage{
		Version:  conn.Version,
		Action:   pbs.Action_Auth,
		BodyData: []byte("ok"),
	})
}

func (h *handler) MessageReceive(conn *Conn, mp pbs.MsgPackage) {
	m := &pbs.MsgItem{}
	if err := msgFormat.Decode(mp.BodyData, m); err != nil {
		return
	}
	m.DeviceID = conn.DeviceID
	m.UserId = conn.UserID
	m.ClientType = conn.ClientType
	_, _ = rpc_client.LogicInit.MessageReceive(context.TODO(), &pbs.MessageReceiveReq{
		Message: m,
	})
}

func (h *handler) MessageACK(conn *Conn, mp pbs.MsgPackage) {
	m := &pbs.MsgItem{}
	if err := msgFormat.Decode(mp.BodyData, m); err != nil {
		return
	}

	h.DeliverMessageByUIDAndDID(m.UserId, m.DeviceID, mp)
}

func (h *handler) SyncTrigger(conn *Conn, mp pbs.MsgPackage) {
	_, _ = rpc_client.LogicInit.MessageSync(context.TODO(), &pbs.MessageSyncReq{
		UserId:    conn.UserID,
		MessageId: string(mp.BodyData),
	})
}

func (h *handler) Headbeat(conn *Conn) {
	conn.Write(pbs.MsgPackage{
		Version: conn.Version,
		Action:  pbs.Action_Headbeat,
	})
}

// 投递消息
func (h *handler) DeliverMessageByUID(UserId int64, mp pbs.MsgPackage) error {
	// 获取设备对应的TCP连接
	conns := TCPServer.LoadConnsByUID(UserId)
	if conns == nil {
		logger.Sugar.Warn("conn id nil")
		return nil
	}
	for _, conn := range conns {
		// 发送消息
		conn.Write(mp)
	}
	return nil
}

// 投递消息
func (h *handler) DeliverMessageByUIDAndDID(UserId int64, DeviceID string, mp pbs.MsgPackage) error {
	// 获取设备对应的TCP连接
	conns := TCPServer.LoadConnsByUID(UserId)
	if conns == nil {
		logger.Sugar.Warn("conn id nil")
		return nil
	}
	for _, conn := range conns {
		if DeviceID == conn.DeviceID {
			// 发送消息
			conn.Write(mp)
		}
	}
	return nil
}

func (h *handler) DeliverMessageByUIDAndNotDID(UserId int64, DeviceID string, mp pbs.MsgPackage) error {
	// 获取设备对应的TCP连接
	conns := TCPServer.LoadConnsByUID(UserId)
	if conns == nil {
		logger.Sugar.Warn("conn id nil")
		return nil
	}
	for _, conn := range conns {
		if conn.DeviceID == DeviceID {
			continue
		}
		// 发送消息
		conn.Write(mp)
	}
	return nil
}
