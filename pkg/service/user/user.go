package user

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type user struct {
	redisPool *redis.Pool
}

var User user

// NewUid 创建一个Uid;len：缓冲池大小()
// db:数据库连接
// businessId：业务id
// len：缓冲池大小(长度可控制缓存中剩下多少id时，去DB中加载)
func InitUser(rdb *redis.Pool) error {
	User = user{
		redisPool: rdb,
	}
	return nil
}

func (s *user) SetUserOnline(uid int64, status bool, deviceType string) error {
	rconn := s.redisPool.Get()
	redis_key := "user_online_status:" + deviceType
	// 用户在线状态
	var statu int
	if status {
		statu = 1
	} else {
		statu = 0
	}
	_, err := rconn.Do("SETBIT", redis_key, strconv.Itoa(int(uid)), statu)
	return err
}

func (s *user) GetUserOnline(uid int64) (bool, error) {
	rconn := s.redisPool.Get()
	redis_key := "user_online_status"
	statu, err := rconn.Do("GETBIT", redis_key, strconv.Itoa(int(uid)))
	var status bool
	if statu.(int) == 0 {
		status = false
	} else {
		status = true
	}
	return status, err
}
