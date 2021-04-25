package dao

import "free-im/internal/http_app/model"

type user struct {
}

var User = new(user)

func (d *user) Get(member_id uint, selects ...string) (model model.UserMember, err error) {
	result := Dao.db.Table(model.TableName()).Select(selects).Find(&model, member_id)
	err = result.Error
	return
}
