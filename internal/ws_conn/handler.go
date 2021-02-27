package ws_conn

import (
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
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

func (h *handler) Auth(ctx *Context, mp pbs.MsgPackage) {

}

func (h *handler) MessageReceive(ctx *Context, mp pbs.MsgPackage) {

}

func (h *handler) MessageACK(ctx *Context, mp pbs.MsgPackage) {

}

func (h *handler) SyncTrigger(ctx *Context, mp pbs.MsgPackage) {

}

func (h *handler) Headbeat(ctx *Context) {

}

// 投递消息
func (h *handler) DeliverMessageByUID(UserId string, mp pbs.MsgPackage) error {

	return nil
}

// 投递消息
func (h *handler) DeliverMessageByUIDAndDID(UserId string, DeviceID string, mp pbs.MsgPackage) error {

	return nil
}

func (h *handler) DeliverMessageByUIDAndNotDID(UserId string, DeviceID string, mp pbs.MsgPackage) error {

	return nil
}
