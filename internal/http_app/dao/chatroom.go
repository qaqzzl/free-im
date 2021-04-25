package dao

import (
	"fmt"
	"free-im/internal/http_app/model"
	"math/rand"
	"strconv"
	"time"
)

type chatroom struct {
}

var Chatroom = new(chatroom)

// * 获取聊天室成员
func (d *chatroom) GetMembers(chatroom_id string) (uids []uint, err error) {
	members, err := Dao.redis.Get().Do("SMEMBERS", "set_im_chatroom_member:"+chatroom_id)
	for _, v := range members.([]interface{}) {
		user_id, _ := strconv.Atoi(string(v.([]uint8)))
		uids = append(uids, uint(user_id))
	}
	return
}

// * 创建群组
func (d *chatroom) CreateGroup(group model.Group) (group_id uint, err error) {
	rand.Seed(time.Now().Unix())
	group.Id = fmt.Sprintf("%06d", rand.Int31n(10000))
	result := Dao.db.Table("`group`").Create(&group)
	group_id = group.GroupId
	err = result.Error
	// redis
	Dao.redis.Get().Do("SADD", "set_im_chatroom_member:"+group.ChatroomId, group.OwnerMemberId) //创建聊天室
	return
}
