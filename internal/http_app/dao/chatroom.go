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
func (d *chatroom) GetMembers(chatroom_id string) (uids []int, err error) {
	members, err := Dao.redis.Get().Do("SMEMBERS", "set_im_chatroom_member:"+chatroom_id)
	for _, v := range members.([]interface{}) {
		user_id, _ := strconv.Atoi(string(v.([]uint8)))
		uids = append(uids, user_id)
	}
	return
}

// * 创建群组
func (d *chatroom) CreateGroup(group model.Group) (group_id string, err error) {
	group_data := make(map[string]string)
	rand.Seed(time.Now().Unix())
	group_data["id"] = fmt.Sprintf("%06d", rand.Int31n(10000))
	group_data["name"] = group.Name
	group_data["avatar"] = group.Avatar
	group_data["chatroom_id"] = group.ChatroomId
	group_data["owner_member_id"] = group.OwnerMemberId
	group_data["founder_member_id"] = group.FounderMemberId
	group_data["created_at"] = strconv.Itoa(int(time.Now().Unix()))
	result, err := Dao.db.Table("`group`").Insert(group_data)
	id, _ := result.LastInsertId()
	group_id = strconv.Itoa(int(id))
	// redis
	Dao.redis.Get().Do("SADD", "set_im_chatroom_member:"+group.ChatroomId, group.OwnerMemberId) //创建聊天室
	return
}
