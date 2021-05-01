package id

import (
	"errors"
	"free-im/pkg/protos/pbs"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type chatroomId struct {
	redisPool  *redis.Pool
	redisKey   string
	errTimeOut error
}

var ChatroomID chatroomId

// NewUid 创建一个Uid;len：缓冲池大小()
// db:数据库连接
// businessId：业务id
// len：缓冲池大小(长度可控制缓存中剩下多少id时，去DB中加载)
func initChatroomID(rdb *redis.Pool) error {
	ChatroomID = chatroomId{
		redisPool:  rdb,
		redisKey:   "chatroom_id",
		errTimeOut: errors.New("get chatroom_id timeout"),
	}
	return nil
}

// 获取聊天室ID
func (id *chatroomId) GetID(sessionType pbs.ChatroomType) (int64, error) {
	rconn := id.redisPool.Get()
	max, err := rconn.Do("get", "id:"+id.redisKey+":max")
	if err != nil {
		return 0, id.errTimeOut
	}
	var maxNum int64
	if max == nil {
		maxNum = 100000
	} else {
		maxNum = byteUintToint64(max.([]uint8))
	}
	//检查号池是否已存在
	exists, err := rconn.Do("exists", "id:"+id.redisKey+":pond")
	if err != nil {
		return 0, id.errTimeOut
	}
	if exists.(int64) == 0 {
		addlen := 1000
		if _, err := rconn.Do("set", "id:"+id.redisKey+":max", maxNum+int64(addlen)); err != nil {
			return 0, id.errTimeOut
		}
		for i := 0; i < addlen; i++ {
			maxNum++
			rconn.Do("sAdd", "id:"+id.redisKey+":pond", maxNum)
		}
	}
	//随机抽取一个号
	cid, err := rconn.Do("sPop", "id:"+id.redisKey+":pond")
	if err != nil {
		return 0, id.errTimeOut
	}
	cidint, _ := strconv.Atoi(strconv.Itoa(int(sessionType)) + byteUintToString(cid.([]uint8)))
	return int64(cidint), nil
}
