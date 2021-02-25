package ws_conn

import (
	"free-im/pkg/protos/pbs"
)

type handler struct{}

var Handler = new(handler)

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
