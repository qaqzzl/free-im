package service

import (
	"errors"
	"fmt"
	"free-im/internal/http_app/dao"
	"free-im/internal/http_app/model"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/util/id"
	"strconv"
)

type ChatRoomService struct {
}

func (s *ChatRoomService) GetChatroomBaseInfo(chatroom_id string, chatroom_type string, user_id string) (resres map[string]string, err error) {
	res := make(map[string]string)
	switch chatroom_type {
	case "1":
		// 查询聊天室成员
		members, _ := dao.Chatroom.GetMembers(chatroom_id)
		for _, v := range members {
			to_user_id := strconv.Itoa(v)
			if to_user_id == user_id {
				continue
			}
			user_member, _ := dao.NewMysql().Table("user_member").
				Where("member_id = " + to_user_id).
				First("nickname,avatar")
			res["name"] = user_member["nickname"]
			res["avatar"] = user_member["avatar"]
		}
	case "2":

	}
	return
}

// 通过好友ID 获取 聊天室ID
func (s *ChatRoomService) FriendIdGetChatroomId(member_id string, friend_id string) (chatroom_id string, err error) {
	var field string
	if member_id > friend_id {
		field = member_id + "," + friend_id
	} else {
		field = friend_id + "," + member_id
	}
	var res interface{}
	if res, err = dao.GetRConn().Do("HGET", "hash_im_chatroom_friend_id_get_chatroom_id", field); err != nil {
		logger.Sugar.Error("获取聊天室ID失败", err)
		return "", err
	}

	if res == nil {
		//生成聊天室ID
		res_chatroom_id, _ := id.ChatroomID.GetID(pbs.ChatroomType_Single)
		chatroom_id = strconv.Itoa(int(res_chatroom_id))
		dao.GetRConn().Do("SADD", "set_im_chatroom_member:"+chatroom_id, member_id, friend_id)      //创建聊天室
		dao.GetRConn().Do("HSET", "hash_im_chatroom_friend_id_get_chatroom_id", field, chatroom_id) //创建单聊跟聊天室关系
	} else {
		chatroom_id = string(res.([]uint8))
	}
	return chatroom_id, err
}

// 聊天室列表
func (s *ChatRoomService) ChatroomList(member_id string) {

}

// 创建群组
func (s *ChatRoomService) CreateGroup(member_id string, group model.Group) (group_id uint, err error) {
	// db
	res_chatroom_id, _ := id.ChatroomID.GetID(pbs.ChatroomType_Group)
	chatroom_id := strconv.Itoa(int(res_chatroom_id))
	group.ChatroomId = chatroom_id
	group.OwnerMemberId = member_id
	group.FounderMemberId = member_id
	group_id, err = dao.Chatroom.CreateGroup(group)
	return group_id, err
}

// 加入群组
func (s *ChatRoomService) AddGroup(member_id string, group_id string, remark string) (ret map[string]string, err error) {
	group, err := dao.NewMysql().Table("`group`").Where(fmt.Sprintf("group_id = %s ", group_id)).First("chatroom_id")
	if len(group) == 0 {
		return ret, errors.New("群组不存在")
	}
	dao.GetRConn().Do("SADD", "set_im_chatroom_member:"+group["chatroom_id"], member_id) //加入聊天室

	ret = make(map[string]string)
	ret["code"] = "0"
	ret["message"] = "加入成功"
	return ret, nil
}

// 我的群组列表
func (s *ChatRoomService) MyGroupList(member_id string) (list []map[string]string, err error) {
	list, err = dao.NewMysql().Table("group_member as gm").Where("gm.member_id = " + member_id + " and gm.status = 'normal'").
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
