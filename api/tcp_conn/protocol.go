package tcp_conn

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"net"
)

func (c *Context) Write(conn net.Conn, p Package) (int, error) {
	// fmt.Println(strings.TrimSpace(string(p.BodyData)))
	// 读取消息的长度
	var length = int32(len(p.BodyData))
	var pkg = new(bytes.Buffer)
	//写入消息头
	if err := binary.Write(pkg, binary.BigEndian, p.Version); err != nil {
		return 0, err
	}
	if err := binary.Write(pkg, binary.BigEndian, p.Action); err != nil {
		return 0, err
	}
	if err := binary.Write(pkg, binary.BigEndian, p.SequenceId); err != nil {
		return 0, err
	}
	if err := binary.Write(pkg, binary.BigEndian, length); err != nil {
		return 0, err
	}

	//写入消息体
	if err := binary.Write(pkg, binary.BigEndian, p.BodyData); err != nil {
		return 0, err
	}
	nn, err := conn.Write(pkg.Bytes())
	if err != nil {
		return 0, err
	}
	return nn, nil
}

func (c *Context) Read() (p Package, err error) {
	// Peek 返回缓存的一个切片，该切片引用缓存中前 n 个字节的数据，
	// 该操作不会将数据读出，只是引用，引用的数据在下一次读取操作之
	// 前是有效的。如果切片长度小于 n，则返回一个错误信息说明原因。
	// 如果 n 大于缓存的总大小，则返回 ErrBufferFull。
	var headInt = 13
	headByte, err := c.r.Peek(headInt)
	if err != nil {
		return p, err
	}
	//创建 Buffer缓冲器
	lengthBuff := bytes.NewBuffer(headByte[9:13])
	var Bodylength int32
	// 通过Read接口可以将buf中得内容填充到data参数表示的数据结构中
	err = binary.Read(lengthBuff, binary.BigEndian, &Bodylength)
	if err != nil {
		return p, err
	}
	// Buffered 返回缓存中未读取的数据的长度

	// debug
	len := Bodylength+int32(headInt)
	if int(len) > math.MaxInt32 || int(len) < 0 {
		len = 0
		// debug
		log.Panicln("数据量超出: ", Bodylength, "  ,  "+ string(p.BodyData))
	}
	// end debug
	if int32(c.r.Buffered()) < len {
		return p, err
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
	n, err := c.r.Read(pack)
	if err != nil {
		return p, err
	}
	Buff := bytes.NewBuffer(headByte[0:4])
	err = binary.Read(Buff, binary.BigEndian, &p.Version)
	Buff = bytes.NewBuffer(headByte[4:5])
	err = binary.Read(Buff, binary.BigEndian, &p.Action)
	Buff = bytes.NewBuffer(headByte[5:9])
	err = binary.Read(Buff, binary.BigEndian, &p.SequenceId)
	p.BodyLength = Bodylength
	p.BodyData = pack[13:n]
	// end debug
	return p, nil
}