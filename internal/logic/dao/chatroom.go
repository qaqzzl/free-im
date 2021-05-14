package dao

import (
	"strconv"
)

type chatroom struct {
}

var Chatroom = new(chatroom)

// * 获取聊天室成员
func (d *chatroom) GetMembers(chatroom_id int64) (member_ids []int64, err error) {
	members, err := Dao.Ris().Do("SMEMBERS", "set_im_chatroom_member:"+strconv.Itoa(int(chatroom_id)))
	for _, v := range members.([]interface{}) {
		user_id, _ := strconv.Atoi(string(v.([]uint8)))
		member_ids = append(member_ids, int64(user_id))
	}
	return
}
