package ws_conn

import (
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
)

type handler struct{}

var Handler = new(handler)

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

func (h *handler) Auth(conn *Conn, mp pbs.MsgPackage) {

}

func (h *handler) MessageReceive(conn *Conn, mp pbs.MsgPackage) {

}

func (h *handler) MessageACK(conn *Conn, mp pbs.MsgPackage) {

}

func (h *handler) SyncTrigger(conn *Conn, mp pbs.MsgPackage) {

}

func (h *handler) Headbeat(conn *Conn) {

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
