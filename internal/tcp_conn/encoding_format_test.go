package tcp_conn

import (
	"fmt"
	"free-im/pkg/protos/pbs"
	"testing"
)

func TestMsgFormat_Encode(t *testing.T) {
	by, err := MsgFormat.Encode(&pbs.MsgACK{
		MessageId: "test",
		UserId:    "1",
		DeviceID:  "2",
	})
	if err != nil {
		t.Error("编码失败", err)
	}
	fmt.Println(by)
}

func TestMsgFormat(t *testing.T) {
	by, err := MsgFormat.Encode(&pbs.MsgACK{
		MessageId: "test",
		UserId:    "1",
		DeviceID:  "2",
	})
	if err != nil {
		t.Error("编码失败", err)
	}
	fmt.Println(by)
	msg := &pbs.MsgACK{}
	err = MsgFormat.Decode(by, msg)
	if err != nil {
		t.Error("解码失败", err)
	}
	fmt.Println(msg)
}
