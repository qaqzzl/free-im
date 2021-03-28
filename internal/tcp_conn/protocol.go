package tcp_conn

import (
	"bytes"
	"encoding/binary"
	"errors"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
)

type protocol struct {
	headerLen int
	bodyLen   int
}

var Protocol = protocol{
	headerLen: 13,
	//bodyLen: 2048,
	bodyLen: 2048 * 100000, // dubug
}

func (p *protocol) Decode(conn *Conn) (mp pbs.MsgPackage, err error) {
	// Peek 返回缓存的一个切片，该切片引用缓存中前 n 个字节的数据，
	// 该操作不会将数据读出，只是引用，引用的数据在下一次读取操作之
	// 前是有效的。如果切片长度小于 n，则返回一个错误信息说明原因。
	// 如果 n 大于缓存的总大小，则返回 ErrBufferFull。
	var headInt = p.headerLen
	headByte, err := conn.r.Peek(headInt)
	if err != nil {
		return mp, err
	}
	//创建 Buffer缓冲器
	lengthBuff := bytes.NewBuffer(headByte[9:13])
	var Bodylength int32
	// 通过Read接口可以将buf中得内容填充到data参数表示的数据结构中
	err = binary.Read(lengthBuff, binary.BigEndian, &Bodylength)
	if err != nil {
		return mp, err
	}

	len := Bodylength + int32(headInt)
	if int(len) > p.bodyLen || int(len) <= 0 {
		// debug
		logger.Sugar.Error("body len error: ", Bodylength)
		return mp, errors.New("body len error")
	}

	// Buffered 返回缓存中未读取的数据的长度
	if int32(conn.r.Buffered()) < len {
		return mp, err
	}
	// 读取消息真正的内容
	pack := make([]byte, int(Bodylength+int32(headInt)))
	// Read 从 b 中读出数据到 p 中，返回读出的字节数和遇到的错误。
	// 如果缓存不为空，则只能读出缓存中的数据，不会从底层 io.Reader
	// 中提取数据，如果缓存为空，则：
	// 1、len(p) >= 缓存大小，则跳过缓存，直接从底层 io.Reader 中读
	// 出到 p 中。
	// 2、len(p) < 缓存大小，则先将数据从底层 io.Reader 中读取到缓存
	// 中，再从缓存读取到 p 中。
	n, err := conn.r.Read(pack)
	if err != nil {
		return mp, err
	}
	var action int8
	Buff := bytes.NewBuffer(headByte[0:4])
	err = binary.Read(Buff, binary.BigEndian, &mp.Version)
	Buff = bytes.NewBuffer(headByte[4:5])
	err = binary.Read(Buff, binary.BigEndian, &action)
	Buff = bytes.NewBuffer(headByte[5:9])
	err = binary.Read(Buff, binary.BigEndian, &mp.SequenceId)

	mp.Action = actionTo(action)
	mp.BodyLength = Bodylength
	mp.BodyData = pack[13:n]
	if mp.BodyData == nil {
		mp.BodyData = []byte("not nil")
	}
	return mp, nil
}

func (p *protocol) Encode(mp pbs.MsgPackage) ([]byte, error) {
	// fmt.Println(strings.TrimSpace(string(p.BodyData)))
	// 读取消息的长度
	var length = int32(len(mp.BodyData))
	var pkg = new(bytes.Buffer)
	action := actionToint8(mp.Action)

	//写入消息头
	if err := binary.Write(pkg, binary.BigEndian, mp.Version); err != nil {
		return nil, err
	}
	if err := binary.Write(pkg, binary.BigEndian, action); err != nil {
		return nil, err
	}
	if err := binary.Write(pkg, binary.BigEndian, mp.SequenceId); err != nil {
		return nil, err
	}
	if err := binary.Write(pkg, binary.BigEndian, length); err != nil {
		return nil, err
	}

	//写入消息体
	if err := binary.Write(pkg, binary.BigEndian, mp.BodyData); err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

func actionToint8(action pbs.Action) int8 {
	var ac int8
	switch action {
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
