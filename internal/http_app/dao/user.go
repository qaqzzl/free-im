package dao

import "free-im/internal/http_app/model"

type user struct {
}

var User = new(user)

func (d *user) GetByUID(member_id uint, selects ...string) (UserMember model.UserMember, err error) {
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
func (d *user) AddFriend() {

}
