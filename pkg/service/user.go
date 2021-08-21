package service

import "strconv"

func (s *service) SetUserOnline(uid int64, status bool) error {
	rconn := s.redisPool.Get()
	redis_key := "user_online_status"
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

func (s *service) GetUserOnline(uid int64) (bool, error) {
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
