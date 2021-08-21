package service

import (
	"errors"
	"free-im/internal/http_app/dao"
	"free-im/internal/http_app/model"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/service/id"
	"gorm.io/gorm"
	"strconv"
)

type ChatRoomService struct {
}

func (s *ChatRoomService) GetChatroomBaseInfo(chatroom_id int64, member_id int64) (res map[string]interface{}, err error) {
	chatroom_id_str := strconv.Itoa(int(chatroom_id))
	chatroom_type := chatroom_id_str[len(chatroom_id_str)-1:]
	res = make(map[string]interface{})
	switch chatroom_type {
	case "1":
		// 查询聊天室成员
		members, _ := dao.Chatroom.GetMembers(chatroom_id)
		for _, v := range members {
			to_member_id := v
			if to_member_id == member_id {
				continue
			}
			user_member, _ := dao.User.GetByUID(to_member_id, "nickname", "avatar", "member_id")

			res["title"] = user_member.Nickname
			res["picture"] = user_member.Avatar
			res["id"] = user_member.MemberId
		}
	case "2":
		group := model.Group{}
		result := dao.Dao.DB().Table(group.TableName()).First(&group)
		if result.Error != nil {
			return nil, result.Error
		}
		res["title"] = group.Name
		res["picture"] = group.Avatar
		res["id"] = group.GroupId
	}
	res["chatroom_type"], _ = strconv.Atoi(chatroom_type)
	return
}

// 通过好友ID 获取 聊天室ID
func (s *ChatRoomService) FriendIdGetChatroomId(member_id int64, friend_id int64) (chatroom_id int64, err error) {
	if member_id == friend_id {
		return
	}
	var field string
	if member_id > friend_id {
		field = strconv.Itoa(int(member_id)) + "," + strconv.Itoa(int(friend_id))
	} else {
		field = strconv.Itoa(int(friend_id)) + "," + strconv.Itoa(int(member_id))
	}
	var res interface{}
	if res, err = dao.Dao.Ris().Do("HGET", "hash_im_chatroom_friend_id_get_chatroom_id", field); err != nil {
		logger.Sugar.Error("获取聊天室ID失败", err)
		return 0, err
	}
	if res == nil {
		//生成聊天室ID
		chatroom_id, err = id.ChatroomID.GetID(pbs.ChatroomType_Single)
		if err == nil {
			dao.Dao.Ris().Do("SADD", "set_im_chatroom_member:"+strconv.Itoa(int(chatroom_id)), member_id, friend_id) //创建聊天室
			dao.Dao.Ris().Do("HSET", "hash_im_chatroom_friend_id_get_chatroom_id", field, chatroom_id)               //创建单聊跟聊天室关系
		}
	} else {
		int_chatroom_id, _ := strconv.Atoi(string(res.([]uint8)))
		chatroom_id = int64(int_chatroom_id)
	}
	return chatroom_id, err
}

// 聊天室列表
func (s *ChatRoomService) ChatroomList(member_id int64) {

}

// 创建群组
func (s *ChatRoomService) CreateGroup(member_id int64, group *model.Group, member_list []int64) (err error) {
	chatroom_id, _ := id.ChatroomID.GetID(pbs.ChatroomType_Group)
	group.ChatroomId = chatroom_id
	group.OwnerMemberId = member_id
	group.FounderMemberId = member_id
	if err = dao.Chatroom.CreateGroup(group); err != nil {
		return
	}
	// 添加群成员
	for _, v := range member_list {
		dao.Chatroom.JoinGroup(&model.GroupMember{
			GroupId:        group.GroupId,
			MemberId:       v,
			MemberIdentity: "common",
			Status:         "normal",
			NotifyLevel:    0,
		}, chatroom_id)
	}

	return
}

// 加入群组
func (s *ChatRoomService) JoinGroup(member_id int64, group_id int64, remark string) (ret map[string]string, err error) {
	ret = make(map[string]string)

	group, err := dao.Chatroom.GetGroupByID(group_id, "*")
	if err != nil {
		return ret, errors.New("系统忙，请稍后再试")
	}
	if group.ChatroomId == 0 {
		return ret, errors.New("群组不存在")
	}
	var GroupMember model.GroupMember
	result := dao.Dao.DB().Where("member_id = ? and group_id = ?", member_id, group_id).First(&GroupMember, "status")
	if result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	if GroupMember.Status == "blacklist" {
		return ret, errors.New("此群组拒绝你的申请🙈")
	} else {
		ret["status"] = "0"
		ret["message"] = "加入成功"
		return
	}
	var status string
	if group.Permissions == "public" {
		status = "normal"
	} else {
		status = "wait"
	}
	// 加入聊天室
	dao.Chatroom.JoinGroup(&model.GroupMember{
		MemberId: member_id, GroupId: group_id,
		MemberIdentity: "common", Status: status,
		NotifyLevel: 0,
	}, group.ChatroomId)

	ret["status"] = "0"
	ret["message"] = "加入成功"
	return
}

// 添加群成员
func (s *ChatRoomService) AddGroupMember(member_id int64, group_id int64, memberList []int64) (groupMembers []model.GroupMember, err error) {
	group, err := dao.Chatroom.GetGroupByID(group_id, "*")
	if err != nil {
		return groupMembers, errors.New("系统忙，请稍后再试")
	}
	if group.ChatroomId == 0 {
		return groupMembers, errors.New("群组不存在")
	}
	for _, v := range memberList {
		var groupMember = model.GroupMember{
			MemberId: v, GroupId: group_id,
			MemberIdentity: "common", Status: "normal",
			NotifyLevel: 0,
		}
		// 加入聊天室
		dao.Chatroom.JoinGroup(&groupMember, group.ChatroomId)
		groupMembers = append(groupMembers, groupMember)
	}
	return
}

// 会员群组列表
func (s *ChatRoomService) MemberGroupList(member_id int64) (Groups []model.Group, err error) {
	Groups, err = dao.Chatroom.MemberGroupListByUID(member_id)
	return
}

// 群组信息
func (s *ChatRoomService) GroupInfo(group_id int64) (Group model.Group, err error) {
	return dao.Chatroom.GroupInfo(group_id)
}

// 搜索群组
func (s *ChatRoomService) SearchGroup(search string) (list []map[string]interface{}, err error) {
	gm := model.Group{}
	//list = make([]map[string]interface{}, 0)
	dao.Dao.DB().Table(gm.TableName()).Where("name like ? ", "%"+search+"%").Find(&list)
	return
}

// 群组成员列表
func (s *ChatRoomService) GroupMember(member_id int64, group_id int64) (list []map[string]interface{}, err error) {
	gm := model.GroupMember{}
	//list = make([]map[string]interface{}, 0)
	dao.Dao.DB().Table(gm.TableName()).Where("group_id = ? and status = ?", group_id, "normal").Find(&list)
	return
}
