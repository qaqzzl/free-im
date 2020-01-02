package service

import (
	"errors"
	"fmt"
	"free-im/dao"
	"free-im/model"
	uuid "github.com/satori/go.uuid"
	"log"
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
		chatroom_id = uuid.NewV4().String()
		rconn.Do("SADD", "set_im_chatroom_member_"+chatroom_id, member_id, friend_id)			//创建聊天室
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
	rconn.Do("SADD", "set_im_chatroom_member_"+chatroom_id, member_id)			//创建聊天室

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
	rconn.Do("SADD", "set_im_chatroom_member_"+group["chatroom_id"], member_id)			//加入聊天室

	ret = make(map[string]string)
	ret["code"] = "0"
	ret["message"] = "加入成功"
	return ret,nil
}