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
