package dao

import (
	"strconv"
)

type chatroom struct {
}

var Chatroom = new(chatroom)

// * 获取聊天室成员
func (d *chatroom) GetMembers(chatroom_id string) (uids []int, err error) {
	members, err := Dao.redis.Get().Do("SMEMBERS", "set_im_chatroom_member:"+chatroom_id)
	for _, v := range members.([]interface{}) {
		user_id, _ := strconv.Atoi(string(v.([]uint8)))
		uids = append(uids, user_id)
	}
	return
}
