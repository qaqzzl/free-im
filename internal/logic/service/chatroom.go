package service

import (
	"errors"
	"fmt"
	"free-im/api/model"
	"free-im/dao"
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type ChatRoomService struct {

}

// 通过好友ID 获取 聊天室ID
func (s *ChatRoomService) FriendIdGetChatroomId(member_id string, friend_id string) (chatroom_id string, err error) {
	rconn := dao.NewRedis()
	var field string
	if member_id > friend_id {
		field = member_id +","+ friend_id
	} else {
		field = friend_id +","+ member_id
	}
	var res interface{}
	if res,err = rconn.Do("HGET", "hash_im_chatroom_friend_id_get_chatroom_id", field); err != nil {
		log.Println(err)
	}

	if res == nil {
		//生成聊天室ID
		chatroom_id = uuid.NewV4().String() + ":ordinary"
		rconn.Do("SADD", "set_im_chatroom_member:"+chatroom_id, member_id, friend_id)			//创建聊天室
		rconn.Do("HSET", "hash_im_chatroom_friend_id_get_chatroom_id", field, chatroom_id)		//创建单聊跟聊天室关系
	} else {
		chatroom_id = string( res.([]uint8) )
	}
	return chatroom_id, err
}


// 聊天室列表
func (s *ChatRoomService) ChatroomList(member_id string) {

}


// 创建群组
func (s *ChatRoomService) CreateGroup(member_id string, group model.Group) (group_id string, err error) {
	// db
	chatroom_id := uuid.NewV4().String()
	group_data := make(map[string]string)
	rand.Seed(time.Now().Unix())
	group_data["id"] = fmt.Sprintf("%06d",rand.Int31n(10000))
	group_data["name"] = group.Name
	group_data["avatar"] = group.Avatar
	group_data["chatroom_id"] = chatroom_id
	group_data["owner_member_id"] = member_id
	group_data["founder_member_id"] = member_id
	group_data["created_at"] = strconv.Itoa(int(time.Now().Unix()))
	result,err := dao.NewMysql().Table("`group`").Insert(group_data)
	id, _ := result.LastInsertId()
	group_id = strconv.Itoa(int(id))

	// redis
	rconn := dao.NewRedis()
	rconn.Do("SADD", "set_im_chatroom_member:"+chatroom_id, member_id)			//创建聊天室

	return group_id, err
}

// 加入群组
func (s *ChatRoomService) AddGroup(member_id string, group_id string, remark string) (ret map[string]string, err error) {
	group,err  := dao.NewMysql().Table("`group`").Where( fmt.Sprintf("group_id = %s ", group_id)).First("chatroom_id")
	if len(group) == 0 {
		return ret,errors.New("群组不存在")
	}
	// redis
	rconn := dao.NewRedis()
	rconn.Do("SADD", "set_im_chatroom_member:"+group["chatroom_id"], member_id)			//加入聊天室

	ret = make(map[string]string)
	ret["code"] = "0"
	ret["message"] = "加入成功"
	return ret,nil
}

// 我的群组列表
func (s *ChatRoomService) MyGroupList(member_id string) (list []map[string]string, err error) {
	list, err = dao.NewMysql().Table("group_member as gm").Where("gm.member_id = "+member_id + " and gm.status = 'normal'").
		Join("INNER JOIN `group` g ON g.group_id=gm.group_id").
		Select("g.group_id,g.name,g.avatar,g.id,g.chatroom_id,g.owner_member_id,g.founder_member_id,g.permissions,g.created_at").
		Get()
	if len(list) == 0 {
		list = make([]map[string]string, 0)
	}
	return list, err
}

// 搜索群组
func (s *ChatRoomService) SearchGroup(search string) (list []map[string]string, err error) {
	return list, err
}