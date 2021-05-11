package dao

import (
	"free-im/internal/http_app/model"
)

type user struct {
}

var User = new(user)

func (d *user) GetByUID(member_id int64, selects ...string) (UserMember model.UserMember, err error) {
	result := Dao.DB().Table(UserMember.TableName()).Select(selects).Find(&UserMember, member_id)
	err = result.Error
	return
}

// * 查询账号是否存在
func (d *user) IsExistAccount(account string, _type string) (b bool, err error) {
	var c int64
	result := Dao.DB().Table("user_auths").Where(`identity_type = ?' and identifier = ?`, account, _type).Limit(1).Count(&c)
	err = result.Error
	if c > 0 {
		return true, err
	}
	return
}

// * 添加好友
func (d *user) AddFriend(member_id int64, friend_id int64, remark string) (err error) {
	var friend model.UserFriendApply
	friend.MemberId = member_id
	friend.FriendId = friend_id
	friend.Remark = remark
	result := Dao.DB().Table("user_friend_apply").Create(&friend)
	err = result.Error
	return err
}

// * 查询好友关系状态
// return int 0:正常好友关系, 1:删除, 2: 非好友关系
func (d *user) QueryFriendBindStatus(member_id int64, to_member_id int64) (int, error) {
	var friend model.UserFriend
	result := Dao.DB().Table("user_friend").Where("member_id = ? and friend_id = ?", member_id, to_member_id).Select("status").First(&friend)
	if result.Error != nil {
		return 2, nil
	}
	return friend.Status, nil
}
