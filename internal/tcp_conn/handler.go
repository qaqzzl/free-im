package tcp_conn

import (
	"context"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/rpc_client"
	"github.com/orcaman/concurrent-map"
)

type handler struct{}

var Handler = new(handler)

func (h *handler) Handler(ctx *Context, mp pbs.MsgPackage) {
	switch mp.Action {
	case pbs.Action_GetMessageID: // 获取消息ID
	case pbs.Action_Auth: // 连接认证
		h.Auth(ctx, mp)
	case pbs.Action_Message: // 消息
		h.MessageReceive(ctx, mp)
	case pbs.Action_MessageACK: // 消息回执
		h.MessageACK(ctx, mp)
	case pbs.Action_SyncTrigger: // 消息同步
		h.SyncTrigger(ctx, mp)
	case pbs.Action_Headbeat: // 心跳
		h.Headbeat(ctx)
	case pbs.Action_Quit:
		ctx.Close()
	default:
		logger.Sugar.Error("Unsupported command:", mp)
	}
}

// client auth handle
func (h *handler) Auth(ctx *Context, mp pbs.MsgPackage) {
	m := pbs.MsgAuth{}
	if err := MsgFormat.Decode(mp.BodyData, &m); err != nil {
		logger.Sugar.Error(err)
		return
	}

	// 伪代码 认证 code ...
	if resp, err := rpc_client.LogicInit.TokenAuth(context.TODO(), &pbs.TokenAuthReq{
		Message: &m,
	}); err != nil {
		logger.Sugar.Error(err)
		return
	} else if resp.Statu == false {
		return
	}

	ctx.IsAuth = true
	ctx.UserID = m.UserID
	ctx.DeviceID = m.DeviceID
	ctx.ClientType = m.ClientType
	ctx.DeviceType = m.DeviceType
	//加入连接集合
	if tmp, ok := TCPServer.ServerConnPool.Get(ctx.UserID); ok {
		device_map := tmp.(cmap.ConcurrentMap)
		// 判断连接是否存在相同设备
		for k, v := range device_map.Items() {
			if k == m.DeviceType { // 如果有同类型的设备登录了 ,通知其设备下线
				vctx := v.(*Context)
				if vctx.DeviceID != m.DeviceID {
					// 通知其设备下线
					msgQuit := &pbs.MsgQuit{
						Title:   "其他设备登陆通知",
						Content: "你的账号在其他设备登陆<br>如不是你本人登陆请<a href=''>修改密码</a>",
					}
					BodyData, _ := MsgFormat.Encode(msgQuit)
					vctx.Write(pbs.MsgPackage{
						Version:  ctx.Version,
						Action:   pbs.Action_Quit,
						BodyData: BodyData,
					})
				}
				// 关闭连接
				vctx.Close()
			}
		}
	}
	TCPServer.StoreConn(ctx)
	ctx.Write(pbs.MsgPackage{
		Version:  ctx.Version,
		Action:   pbs.Action_Auth,
		BodyData: []byte("ok"),
	})
}

func (h *handler) MessageReceive(ctx *Context, mp pbs.MsgPackage) {
	m := &pbs.MsgItem{}
	if err := MsgFormat.Decode(mp.BodyData, m); err != nil {
		return
	}
	m.DeviceID = ctx.DeviceID
	m.UserId = ctx.UserID
	m.ClientType = ctx.ClientType
	_, _ = rpc_client.LogicInit.MessageReceive(context.TODO(), &pbs.MessageReceiveReq{
		Message: m,
	})
}

func (h *handler) MessageACK(ctx *Context, mp pbs.MsgPackage) {
	m := &pbs.MsgItem{}
	if err := MsgFormat.Decode(mp.BodyData, m); err != nil {
		return
	}

	h.DeliverMessageByUIDAndDID(m.UserId, m.DeviceID, mp)
}

func (h *handler) SyncTrigger(ctx *Context, mp pbs.MsgPackage) {
	_, _ = rpc_client.LogicInit.MessageSync(context.TODO(), &pbs.MessageSyncReq{
		UserId:    ctx.UserID,
		MessageId: string(mp.BodyData),
	})
}

func (h *handler) Headbeat(ctx *Context) {
	ctx.Write(pbs.MsgPackage{
		Version: ctx.Version,
		Action:  pbs.Action_Headbeat,
	})
}

// 投递消息
func (h *handler) DeliverMessageByUID(UserId string, mp pbs.MsgPackage) error {
	// 获取设备对应的TCP连接
	ctxconns := TCPServer.LoadConnsByUID(UserId)
	if ctxconns == nil {
		logger.Sugar.Warn("ctx id nil")
		return nil
	}
	for _, ctxconn := range ctxconns {
		// 发送消息
		ctxconn.Write(mp)
	}
	return nil
}

// 投递消息
func (h *handler) DeliverMessageByUIDAndDID(UserId string, DeviceID string, mp pbs.MsgPackage) error {
	// 获取设备对应的TCP连接
	ctxconns := TCPServer.LoadConnsByUID(UserId)
	if ctxconns == nil {
		logger.Sugar.Warn("ctx id nil")
		return nil
	}
	for _, ctxconn := range ctxconns {
		if DeviceID == ctxconn.DeviceID {
			// 发送消息
			ctxconn.Write(mp)
		}
	}
	return nil
}

func (h *handler) DeliverMessageByUIDAndNotDID(UserId string, DeviceID string, mp pbs.MsgPackage) error {
	// 获取设备对应的TCP连接
	ctxconns := TCPServer.LoadConnsByUID(UserId)
	if ctxconns == nil {
		logger.Sugar.Warn("ctx id nil")
		return nil
	}
	for _, ctxconn := range ctxconns {
		if ctxconn.DeviceID == DeviceID {
			continue
		}
		// 发送消息
		ctxconn.Write(mp)
	}
	return nil
}
