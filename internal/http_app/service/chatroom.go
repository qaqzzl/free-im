package service

import (
	"errors"
	"free-im/internal/http_app/dao"
	"free-im/internal/http_app/model"
	"free-im/pkg/id"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"strconv"
)

type ChatRoomService struct {
}

func (s *ChatRoomService) GetChatroomBaseInfo(chatroom_id string, chatroom_type string, member_id int64) (res map[string]string, err error) {
	res = make(map[string]string)
	switch chatroom_type {
	case "1":
		// 查询聊天室成员
		members, _ := dao.Chatroom.GetMembers(chatroom_id)
		for _, v := range members {
			to_member_id := v
			if to_member_id == member_id {
				continue
			}
			user_member, _ := dao.User.GetByUID(to_member_id, "nickname", "avatar")

			res["name"] = user_member.Nickname
			res["avatar"] = user_member.Avatar
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
	if res, err = dao.Dao.Ris().Do("HGET", "hash_im_chatroom_friend_id_get_chatroom_id", field); err != nil {
		logger.Sugar.Error("获取聊天室ID失败", err)
		return "", err
	}

	if res == nil {
		//生成聊天室ID
		res_chatroom_id, _ := id.ChatroomID.GetID(pbs.ChatroomType_Single)
		chatroom_id = strconv.Itoa(int(res_chatroom_id))
		dao.Dao.Ris().Do("SADD", "set_im_chatroom_member:"+chatroom_id, member_id, friend_id)      //创建聊天室
		dao.Dao.Ris().Do("HSET", "hash_im_chatroom_friend_id_get_chatroom_id", field, chatroom_id) //创建单聊跟聊天室关系
	} else {
		chatroom_id = string(res.([]uint8))
	}
	return chatroom_id, err
}

// 聊天室列表
func (s *ChatRoomService) ChatroomList(member_id string) {

}

// 创建群组
func (s *ChatRoomService) CreateGroup(member_id int64, group model.Group) (group_id int64, err error) {
	res_chatroom_id, _ := id.ChatroomID.GetID(pbs.ChatroomType_Group)
	chatroom_id := strconv.Itoa(int(res_chatroom_id))
	group.ChatroomId = chatroom_id
	group.OwnerMemberId = member_id
	group.FounderMemberId = member_id
	group_id, err = dao.Chatroom.CreateGroup(group)
	return group_id, err
}

// 加入群组
func (s *ChatRoomService) JoinGroup(member_id int64, id string, remark string) (ret map[string]string, err error) {
	group, err := dao.Chatroom.GetGroupByID(id, "chatroom_id")
	if err != nil {
		return ret, errors.New("系统忙，请稍后再试")
	}
	if group.ChatroomId == "" {
		return ret, errors.New("群组不存在")
	}
	// 加入聊天室
	dao.Chatroom.JoinGroup(group.ChatroomId, member_id)

	ret = make(map[string]string)
	ret["code"] = "0"
	ret["message"] = "加入成功"
	return
}

// 会员群组列表
func (s *ChatRoomService) MemberGroupList(member_id uint) (MemberGroups []*model.GroupMember, err error) {
	MemberGroups, err = dao.Chatroom.MemberGroupListByUID(member_id)
	return
}

// 搜索群组
func (s *ChatRoomService) SearchGroup(search string) (list []map[string]string, err error) {
	return
}
