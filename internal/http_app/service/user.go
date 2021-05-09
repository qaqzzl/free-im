package service

import (
	"errors"
	"fmt"
	"free-im/internal/http_app/dao"
	"free-im/internal/http_app/model"
	"time"
)

type UserService struct {
}

// 获取会员信息
func (s *UserService) GetMemberInfo(member_id string) (user_member model.UserMember, err error) {
	result := dao.Dao.DB().Table("user_member").
		Where("member_id = ?", member_id).
		Select("member_id,nickname,gender,birthdate,avatar,signature,city,province").
		Find(&user_member)
	err = result.Error
	return
}

// 修改会员信息
func (s *UserService) UpdateMemberInfo(member_id int64, data model.UserMember) (err error) {
	// 判断用户昵称是否存在
	if data.Nickname == "" {
		return errors.New("昵称不能为空")
	}
	if s.IsMemberNickname(member_id, data.Nickname) {
		return errors.New("昵称已经被使用")
	}
	result := dao.Dao.DB().Table("user_member").
		Where("member_id = ?", member_id).Save(data)
	err = result.Error
	return
}

// 判断用户昵称是否存在
func (s *UserService) IsMemberNickname(member_id int64, nickname string) bool {
	var count int64
	dao.Dao.DB().Table("user_member").
		Where("member_id != ? and nickname = ?", member_id, nickname).Count(&count)
	if count > 0 {
		return true
	}
	return false
}

// 添加好友
func (s *UserService) AddFriend(member_id int64, friend_id int64, remark string) (ret map[string]string, err error) {
	ret = make(map[string]string)
	// 判断是否已经申请
	var user_friend map[string]interface{}
	dao.Dao.DB().Table("user_friend_apply").
		Where("member_id = ? and friend_id = ?", member_id, friend_id).
		Order("id desc").
		Select("member_id", "friend_id", "status").
		Find(&user_friend)
	if len(user_friend) == 0 {
		if err = dao.User.AddFriend(member_id, friend_id, remark); err != nil {
			return ret, err
		}
		ret["message"] = "已发送, 等待对方同意"
		ret["code"] = "1"
		return ret, err
	}

	switch user_friend["status"] {
	case "0":
		ret["message"] = "已发送, 等待对方同意"
		ret["code"] = "1"
		break
	case "1":
		// 判断是否为正常好友关系
		me_is_friend, _ := dao.User.QueryFriendBindStatus(member_id, user_friend["friend_id"].(int64))
		to_is_friend, _ := dao.User.QueryFriendBindStatus(member_id, user_friend["friend_id"].(int64))
		if me_is_friend != 0 || to_is_friend != 0 {
			dao.User.AddFriend(member_id, friend_id, remark)
			ret["message"] = "已发送, 等待对方同意"
			ret["code"] = "1"
		} else {
			ret["message"] = "添加成功, 开始聊天吧"
			ret["code"] = "0"
		}
		break
	case "2":
		if err = dao.User.AddFriend(member_id, friend_id, remark); err != nil {
			return ret, err
		}
		ret["message"] = "已发送, 等待对方同意"
		ret["code"] = "1"
	}
	return ret, err
}

// 删除好友
func (s *UserService) DelFriend(member_id int64, friend_id int64) (err error) {
	// 判断是否已经存在
	result := dao.Dao.DB().Table("user_friend").
		Where("member_id = %s and friend_id = ? or member_id = ? and friend_id = ?", member_id, friend_id, friend_id, member_id).
		Delete(model.UserFriend{})
	err = result.Error
	// 删除聊天室 && 聊天室消息
	return err
}

// 好友申请列表
func (s *UserService) FriendApplyList(member_id int64) (list []map[string]string, err error) {
	result := dao.Dao.DB().Table("user_friend_apply as ufa").Where("ufa.friend_id = ?", member_id).
		Joins("INNER JOIN user_member um ON um.member_id=ufa.member_id").
		Select("um.member_id,um.nickname,um.avatar,um.signature,um.gender,ufa.remark, ufa.status, ufa.id").
		Order("id desc").
		Find(&list)
	if len(list) == 0 {
		list = make([]map[string]string, 0)
	}
	err = result.Error
	return
}

// 好友申请同意/拒绝操作
func (s *UserService) FriendApplyAction(id string, member_id string, action string) (ret map[string]string, err error) {
	ret = make(map[string]string)
	friend_apply, err := dao.NewMysql().Table("user_friend_apply").Where("id = " + id + " and status = 0").
		First("member_id,friend_id")
	if len(friend_apply) == 0 {
		ret["message"] = "操作失败"
		ret["code"] = "1"
		return ret, err
	}
	if friend_apply["friend_id"] != member_id {
		ret["message"] = "违法操作"
		ret["code"] = "1"
		return ret, err
	}
	if _, err = dao.NewMysql().Table("user_friend_apply").Where("id = " + id + " and status = 0").Update("status=" + action); err != nil {
		return ret, err
	}
	if action == "1" {
		timeUnix := time.Now().Unix()
		sql := fmt.Sprintf("INSERT INTO `user_friend` (member_id,friend_id,status,created_at) VALUES (%s,%s,%d,%d) "+
			"ON DUPLICATE KEY UPDATE status=VALUES(status)",
			friend_apply["member_id"], friend_apply["friend_id"], 0, timeUnix)
		if _, err = dao.MysqlConn.Exec(sql); err != nil {
			return ret, err
		}
		sql = fmt.Sprintf("INSERT INTO `user_friend` (member_id,friend_id,status,created_at) VALUES (%s,%s,%d,%d) "+
			"ON DUPLICATE KEY UPDATE status=VALUES(status)",
			friend_apply["friend_id"], friend_apply["member_id"], 0, timeUnix)
		if _, err = dao.MysqlConn.Exec(sql); err != nil {
			return ret, err
		}
	}
	ret["message"] = "操作成功"
	ret["code"] = "0"
	return ret, err
}

// 好友列表
func (s *UserService) FriendList(member_id int64) (list []map[string]string, err error) {
	err = dao.Dao.DB().Table("user_friend").
		Where("member_id = ? and status = 0", member_id).
		Select("member_id,friend_id").
		Find(&list).Error
	if len(list) == 0 {
		list = make([]map[string]string, 0)
	}
	for k, vo := range list {
		var (
			where string
		)
		where = "member_id = " + vo["friend_id"]
		member, _ := dao.NewMysql().Table("user_member").Where(where).First("avatar,gender,member_id,nickname,signature")
		list[k] = member
	}
	return list, err
}

// 搜索好友
func (s *UserService) SearchMember(search string) (list []map[string]string, err error) {
	list, err = dao.NewMysql().Table("user_member").Where("nickname like '%" + search + "%' or id = '" + search + "'").
		Select("member_id,nickname,avatar,signature,gender,birthdate").
		Get()

	if len(list) == 0 {
		list = make([]map[string]string, 0)
	}
	return list, err
}

// 他人基本信息(他人主页)
func (s *UserService) OthersHomeInfo(member_id int64, to_member_id int64) (info map[string]interface{}, err error) {
	user_member := model.UserMember{}
	result := dao.Dao.DB().Table("user_member").Where("member_id = ?", to_member_id).Select("member_id,nickname,gender,birthdate,avatar,signature,city,province").First(&user_member)
	if result.Error != nil {
		return nil, result.Error
	}
	// 判断是否是好友关系
	info["is_friend"] = "no"
	is_friend, _ := dao.User.QueryFriendBindStatus(member_id, to_member_id)
	if is_friend == 0 {
		info["is_friend"] = "yes"
	}
	info["member_info"] = user_member
	return info, err
}
