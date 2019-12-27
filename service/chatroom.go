package service

import (
	"free-im/dao"
	uuid "github.com/satori/go.uuid"
	"log"
)

type ChatRoomService struct {

}

// 通过好友ID 获取 聊天室ID
func (s *ChatRoomService) FriendIdGetChatroomId(member_id string, friend_id string) (chatroom_id string, err error) {
	rconn := dao.Newredis()
	var field string
	if member_id > friend_id {
		field = member_id +","+ friend_id
	} else {
		field = friend_id +","+ member_id
	}
	var res interface{}
	if res,err = rconn.Do("HGET", "hash_im_chatroom_friend_id_get_chatroom_id", field); err == nil {
		log.Println(err)
	}

	if res == nil {
		//生成聊天室ID
		chatroom_id = uuid.NewV4().String()
		rconn.Do("SADD", "set_im_chatroom_member_"+chatroom_id, member_id, friend_id)			//创建聊天室
		rconn.Do("HSET", "hash_im_chatroom_friend_id_get_chatroom_id", field, chatroom_id)						//创建聊天室
	} else {
		chatroom_id = string( res.([]uint8) )
	}
	return chatroom_id, err
}


// 聊天室列表
func (s *ChatRoomService) ChatroomList(member_id string) {

}