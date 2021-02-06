package tcp_conn

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
)

type msgFormat struct {
	Coding string // json, protobuf
}

var MsgFormat = msgFormat{
	Coding: "json",
}

func (f *msgFormat) Decode(data []byte, msg proto.Message) error {
	var err error
	if f.Coding == "json" {
		err = json.Unmarshal(data, &msg)
	} else if f.Coding == "protobuf" {
		err = proto.Unmarshal(data, msg)
	}
	return err
}

func (f *msgFormat) Encode(msg proto.Message) (by []byte, err error) {
	if f.Coding == "json" {
		by, err = json.Marshal(msg)
	} else if f.Coding == "protobuf" {
		by, err = proto.Marshal(msg.(proto.Message))
	}
	return by, err
}
