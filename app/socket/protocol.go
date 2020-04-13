package socket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

func (c *Context) Write(conn net.Conn, message []byte) (int, error) {
	fmt.Println(strings.TrimSpace(string(message)))
	// 读取消息的长度
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	//写入消息头
	err := binary.Write(pkg, binary.BigEndian, length)
	if err != nil {
		return 0, err
	}
	//写入消息体
	err = binary.Write(pkg, binary.BigEndian, message)
	if err != nil {
		return 0, err
	}
	nn, err := conn.Write(pkg.Bytes())
	if err != nil {
		return 0, err
	}
	return nn, nil
}

func (c *Context) Read() ([]byte, error) {
	// Peek 返回缓存的一个切片，该切片引用缓存中前 n 个字节的数据，
	// 该操作不会将数据读出，只是引用，引用的数据在下一次读取操作之
	// 前是有效的。如果切片长度小于 n，则返回一个错误信息说明原因。
	// 如果 n 大于缓存的总大小，则返回 ErrBufferFull。
	lengthByte, err := c.r.Peek(4)
	if err != nil {
		return nil, err
	}
	//创建 Buffer缓冲器
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	// 通过Read接口可以将buf中得内容填充到data参数表示的数据结构中
	err = binary.Read(lengthBuff, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	// Buffered 返回缓存中未读取的数据的长度
	if int32(c.r.Buffered()) < length+4 {
		return nil, err
	}
	// 读取消息真正的内容
	pack := make([]byte, int(4+length))
	// Read 从 b 中读出数据到 p 中，返回读出的字节数和遇到的错误。
	// 如果缓存不为空，则只能读出缓存中的数据，不会从底层 io.Reader
	// 中提取数据，如果缓存为空，则：
	// 1、len(p) >= 缓存大小，则跳过缓存，直接从底层 io.Reader 中读
	// 出到 p 中。
	// 2、len(p) < 缓存大小，则先将数据从底层 io.Reader 中读取到缓存
	// 中，再从缓存读取到 p 中。
	n, err := c.r.Read(pack)
	if err != nil {
		return nil, err
	}
	return pack[4:n], nil
}
