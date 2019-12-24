package service

import (
	"fmt"
	"free-im/dao"
)

type UserService struct {

}

// 获取会员信息
func (s *UserService) GetMemberInfo(member_id string) (user_member map[string]string) {
	return dao.NewMysql().Table("user_member").First("member_id,nickname,gender,birthdate,avatar,signature,city,province")
}

// 添加好友
func (s *UserService) AddFriend(member_id string, friend_id string) (ret map[string]string, err error) {
	ret = make(map[string]string)
	// 判断是否已经申请
	user_friend := dao.NewMysql().Table("user_friend").
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
		if _, err := dao.NewMysql().Table("user_friend").Where(fmt.Sprintf("friend_id = %s and member_id = %s",member_id, friend_id) ).Update("status","1"); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	ret["message"] = "添加成功, 开始聊天吧"
	ret["code"] = "0"
	return ret, err
}