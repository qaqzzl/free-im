package service

import (
	"fmt"
	"free-im/dao"
)

type UserService struct {

}

// 获取会员信息
func (s *UserService) GetMemberInfo(member_id string) (user_member map[string]string, err error) {
	return dao.NewMysql().Table("user_member").First("member_id,nickname,gender,birthdate,avatar,signature,city,province")
}

// 添加好友
func (s *UserService) AddFriend(member_id string, friend_id string) (ret map[string]string, err error) {
	ret = make(map[string]string)
	// 判断是否已经申请
	user_friend,_ := dao.NewMysql().Table("user_friend").
		Where( fmt.Sprintf("member_id = %s and friend_id = %s or member_id = %s and friend_id = %s",member_id, friend_id, friend_id, member_id) ).
		First("member_id,friend_id,status")
	if len(user_friend) == 0 {
		user_friend["member_id"] = member_id
		user_friend["friend_id"] = friend_id
		if _,err = dao.NewMysql().Table("user_friend").Insert(user_friend); err != nil {
			return ret, err
		}
		ret["message"] = "添加成功, 等待对方同意"
		ret["code"] = "1"
		return ret, err
	}

	if user_friend["status"] == "0" {
		if user_friend["friend_id"] != member_id {
			ret["message"] = "添加成功, 等待对方同意"
			ret["code"] = "1"
			return ret, err
		}
		// 如果好友ID是自己, 则直接添加成功
		if _, err := dao.NewMysql().Table("user_friend").Where(fmt.Sprintf("friend_id = %s and member_id = %s",member_id, friend_id) ).Update("status = 1"); err != nil {
			return nil, err
		}
	}

	ret["message"] = "添加成功, 开始聊天吧"
	ret["code"] = "0"
	return ret, err
}

// 删除好友
func (s *UserService) DelFriend(member_id string, friend_id string) (err error) {
	// 判断是否已经存在
	_,err = dao.NewMysql().Table("user_friend").
		Where( fmt.Sprintf("member_id = %s and friend_id = %s or member_id = %s and friend_id = %s",member_id, friend_id, friend_id, member_id) ).
		Delete()
	// 删除聊天室 && 聊天室消息
	return err
}

// 好友申请列表
func (s *UserService) FriendApplyList(member_id string) ([]map[string]string, error) {
	list, err := dao.NewMysql().Table("user_friend as uf").Where("uf.friend_id = "+member_id + " and status = 0").
		Join("INNER JOIN user_member um ON um.member_id=uf.member_id").
		Select("um.member_id,um.nickname,um.avatar,um.signature,um.gender").
		Get()
	if len(list) == 0 {
		list = make([]map[string]string, 0)
	}
	return list, err
}

// 好友列表
func (s *UserService) FriendList(member_id string) (list []map[string]string, err error) {
	list, err = dao.NewMysql().Table("user_friend as uf").Where("uf.friend_id = "+member_id + " and status = 1 or uf.member_id = "+member_id + " and status = 1").
		Join("INNER JOIN user_member um ON um.member_id=uf.member_id").
		Select("um.member_id,um.nickname,um.avatar,um.signature,um.gender").
		Get()
	if len(list) == 0 {
		list = make([]map[string]string, 0)
	}
	return list, err
}
