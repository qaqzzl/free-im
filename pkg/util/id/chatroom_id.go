package id

import (
	"errors"
	"free-im/pkg/protos/pbs"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type logger interface {
	Error(error)
}

// Logger Log接口，如果设置了Logger，就使用Logger打印日志，如果没有设置，就使用内置库log打印日志
var Logger logger

// ErrTimeOut 获取uid超时错误
var ErrTimeOut = errors.New("get uid timeout")

type chatroomId struct {
	rdb redis.Conn
}

var ChatroomID chatroomId

// NewUid 创建一个Uid;len：缓冲池大小()
// db:数据库连接
// businessId：业务id
// len：缓冲池大小(长度可控制缓存中剩下多少id时，去DB中加载)
func InitChatroomID(rdb redis.Conn) error {
	ChatroomID = chatroomId{
		rdb: rdb,
	}
	return nil
}

// 获取聊天室ID
func (id *chatroomId) GetID(sessionType pbs.ChatroomType) (int64, error) {
	max, err := id.rdb.Do("get", "id:chatroom_id:max")
	if err != nil {
		return 0, ErrTimeOut
	}
	var maxNum int64
	if max == nil {
		maxNum = 100000
	} else {
		maxNum = byteUintToint64(max.([]uint8))
	}
	//检查号池是否已存在
	exists, err := id.rdb.Do("exists", "id:chatroom_id:pond")
	if err != nil {
		return 0, ErrTimeOut
	}
	if exists.(int64) == 0 {
		addlen := 1000
		if _, err := id.rdb.Do("set", "id:chatroom_id:max", maxNum+int64(addlen)); err != nil {
			return 0, ErrTimeOut
		}
		for i := 0; i < addlen; i++ {
			maxNum++
			id.rdb.Do("sAdd", "id:chatroom_id:pond", maxNum)
		}
	}
	//随机抽取一个号
	cid, err := id.rdb.Do("sPop", "id:chatroom_id:pond")
	if err != nil {
		return 0, ErrTimeOut
	}
	cidint, _ := strconv.Atoi(strconv.Itoa(int(sessionType)) + byteUintToString(cid.([]uint8)))
	return int64(cidint), nil
}

func byteUintToint64(bs []uint8) int64 {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	in, _ := strconv.Atoi(string(ba))
	return int64(in)
}

func byteUintToString(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}
