package user

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type user struct {
	redisPool *redis.Pool
}

var User user

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
	redis_key := "user_online_status:mobile"
	statuMobile, err := rconn.Do("GETBIT", redis_key, strconv.Itoa(int(uid)))
	if err != nil {
		return false, err
	}
	redis_key = "user_online_status:pc"
	statuPc, err := rconn.Do("GETBIT", redis_key, strconv.Itoa(int(uid)))
	if err != nil {
		return false, err
	}
	var status bool
	if statuMobile.(int64) == 0 && statuPc.(int64) == 0 {
		status = false
	} else {
		status = true
	}
	return status, err
}
