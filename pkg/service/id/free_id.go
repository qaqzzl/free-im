package id

import (
	"errors"
	"free-im/pkg/util"
	"github.com/gomodule/redigo/redis"
)

type freeID struct {
	redisPool  *redis.Pool
	redisKey   string
	errTimeOut error
}

var FreeID freeID

// NewUid 创建一个Uid;len：缓冲池大小()
// db:数据库连接
// businessId：业务id
// len：缓冲池大小(长度可控制缓存中剩下多少id时，去DB中加载)
func InitFreeID(rdb *redis.Pool) error {
	FreeID = freeID{
		redisPool:  rdb,
		redisKey:   "free_id",
		errTimeOut: errors.New("get free_id timeout"),
	}
	return nil
}

// 获取应用号（QQ号）
func (id *freeID) GetID() (string, error) {
	rconn := id.redisPool.Get()
	//检查号池是否已存在
	exists, err := rconn.Do("exists", "id:"+id.redisKey+":pond")
	if err != nil {
		return "", id.errTimeOut
	}
	if exists.(int64) == 0 {
		max, err := rconn.Do("get", "id:"+id.redisKey+":max")
		if err != nil {
			return "", id.errTimeOut
		}
		var maxNum int64
		if max == nil {
			maxNum = 100000
		} else {
			maxNum = util.ByteUintToint64(max.([]uint8))
		}
		addlen := 1000
		if _, err := rconn.Do("set", "id:"+id.redisKey+":max", maxNum+int64(addlen)); err != nil {
			return "", id.errTimeOut
		}
		for i := 0; i < addlen; i++ {
			maxNum++
			rconn.Do("sAdd", "id:"+id.redisKey+":pond", maxNum)
		}
	}
	//随机抽取一个号
	cid, err := rconn.Do("sPop", "id:"+id.redisKey+":pond")
	if err != nil {
		return "", id.errTimeOut
	}
	cidint := util.ByteUintToString(cid.([]uint8))
	return cidint, nil
}
