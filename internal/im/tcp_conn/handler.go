package tcp_conn

import (
	"context"
	"encoding/json"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/rpc_client"
	"github.com/orcaman/concurrent-map"
)

type handler struct{}

var Handler = new(handler)

func (h *handler) Handler(ctx *Context, mp pbs.MessagePackage) {
	switch mp.Action {
	case pbs.Action_GetMessageID: // 获取消息ID
	case pbs.Action_Auth: // 连接认证
		h.Auth(ctx, mp)
	case pbs.Action_Message: // 消息
		h.MessageReceive(ctx, mp)
	case pbs.Action_MessageACK: // 消息回执
	case pbs.Action_SyncTrigger: // 消息同步
	case pbs.Action_Headbeat: // 心跳
		h.Headbeat(ctx)
	case pbs.Action_Quit:
		ctx.TcpConn.Close()
	default:
		logger.Logger.Info("Unsupported command")
	}
}

// client auth handle
func (h *handler) Auth(ctx *Context, mp pbs.MessagePackage) {
	m := pbs.AuthMessage{}
	if err := json.Unmarshal(mp.BodyData, &m); err != nil {
		logger.Sugar.Error(err)
		return
	}
	// 伪代码 认证 code ...
	resp, _ := rpc_client.LogicInit.TokenAuth(context.TODO(), &pbs.TokenAuthReq{
		Message: &m,
	})
	if resp.Statu == false {
		return
	}

	ctx.IsAuth = true
	ctx.UserID = m.UserID
	ctx.DeviceID = m.DeviceID
	ctx.ClientType = m.ClientType
	ctx.DeviceType = m.DeviceType
	clientDevice := ClientDevice{
		DeviceID:   m.DeviceID,
		ClientType: m.ClientType,
		Context:    ctx,
	}
	//加入连接集合
	if tmp, ok := SocketConnPool.Get(ctx.UserID); ok {
		device_map := tmp.(cmap.ConcurrentMap)
		// 判断连接是否存在相同设备
		for k, v := range device_map.Items() {
			if k == m.DeviceType { // 如果有同类型的设备登录了 ,通知其设备下线
				device := v.(ClientDevice)
				if device.DeviceID != m.DeviceID {
					// 通知其设备下线 code...

				}
				// 关闭连接
				device.Context.TcpConn.Close()
				device.Context.IsConnStatus = false // 伪代码
				for {                               // 等待上个连接完全关闭
					if device.Context.IsConnStatus == false {
						break
					}
				}
			}
		}
		device_map.Set(m.DeviceType, clientDevice)
		SocketConnPool.Set(ctx.UserID, device_map)
	} else {
		device_map := cmap.New()
		device_map.Set(m.DeviceType, clientDevice)
		SocketConnPool.Set(ctx.UserID, device_map)
	}
	// 认证成功通知 code ...
	ctx.SendMessage(ctx.TcpConn, pbs.MessagePackage{
		Version: ctx.Version,
		Action:  pbs.Action_Auth,
		BodyData: []byte("ok"),
	})
}

func (h *handler) MessageReceive(ctx *Context, mp pbs.MessagePackage) {
	m := pbs.MessageItem{}
	if err := json.Unmarshal(mp.BodyData, &m); err != nil {
		logger.Sugar.Error(err)
		return
	}
	m.DeviceID = ctx.DeviceID
	m.UserId = ctx.UserID
	m.ClientType = ctx.ClientType
	_, _ = rpc_client.LogicInit.MessageReceive(context.TODO(), &pbs.MessageReceiveReq{
		Message: &m,
	})
}

func (h *handler) Headbeat(ctx *Context) {
	ctx.SendMessage(ctx.TcpConn, pbs.MessagePackage{
		Version: ctx.Version,
		Action:  pbs.Action_Headbeat,
	})
}
