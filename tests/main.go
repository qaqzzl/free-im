package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"free-im/internal/tcp_conn"
	"free-im/pkg/protos/pbs"
	"log"
	"net"
	"time"
)

//发送信息
func sender(conn net.Conn) {
	c := time.Tick(time.Second * 8)
	var SequenceId int32
	SequenceId = 1
	go func() {
		for {
			<-c
			words := pbs.MsgPackage{
				Version:    1,
				Action:     100,
				SequenceId: SequenceId,
				BodyData:   []byte(""),
			}
			if b, err := tcp_conn.Protocol.Encode(words); err == nil {
				conn.Write(b)
				SequenceId++
			}
		}
	}()

	fmt.Println("send over")

	//接收服务端消息
	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), "waiting server back msg error: ", err)
			return
		}
		mp := pbs.MsgPackage{}
		var action int8
		Buff := bytes.NewBuffer(buffer[0:4])
		err = binary.Read(Buff, binary.BigEndian, &mp.Version)
		Buff = bytes.NewBuffer(buffer[4:5])
		err = binary.Read(Buff, binary.BigEndian, &action)
		Buff = bytes.NewBuffer(buffer[5:9])
		err = binary.Read(Buff, binary.BigEndian, &mp.SequenceId)

		mp.Action = actionTo(action)
		mp.BodyData = buffer[13:n]
		// Log(conn.RemoteAddr().String(), "receive server back msg: ", mp)
	}

}

//日志
func Log(v ...interface{}) {
	log.Println(v...)
}

//func NewServer() {
//	server := "101.132.107.212:1208"
//	//server := "127.0.0.1:1208"
//	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
//	if err != nil {
//		Log(os.Stderr, "Fatal error:", err.Error())
//		//os.Exit(1)
//		return
//	}
//	conn, err := net.DialTCP("tcp", nil, tcpAddr)
//	if err != nil {
//		Log(os.Stderr, "Fatal error:", err.Error())
//		//os.Exit(1)
//		return
//	}
//
//	fmt.Println("connection success")
//	go sender(conn)
//}

func main() {
	for {
		NewTicker := time.NewTicker(time.Millisecond * 1)
		<-NewTicker.C
		NewTicker.Stop()
	}

	//for i := 0; i < 800; i++ {
	//	fmt.Println(i)
	//	NewServer()
	//}
	c := time.Tick(time.Second * 100)
	for {
		<-c
	}
}

func actionTo(i int8) (ac pbs.Action) {
	switch i {
	case 0:
		ac = 0
	case 1:
		ac = 1
	case 2:
		ac = 2
	case 3:
		ac = 3
	case 4:
		ac = 4
	case 10:
		ac = 10
	case 11:
		ac = 11
	case 100:
		ac = 100
	}
	return ac
}
