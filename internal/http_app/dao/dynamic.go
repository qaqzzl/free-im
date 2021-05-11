package dao

import "free-im/internal/http_app/model"

type dynamic struct {
}

var Dynamic = new(dynamic)

func (d *dynamic) Create(m model.Dynamic) (dynamic_id int64, err error) {
	result := Dao.DB().Create(&m)
	err = result.Error
	dynamic_id = m.DynamicId
	return
}

func (d *dynamic) DynamicList(page int, prepage int) (total int64, lists []map[string]interface{}) {
	lists = make([]map[string]interface{}, 0)
	Dao.DB().Table("dynamic").Count(&total)
	Dao.DB().Table("dynamic as d").Joins("join user_member um on um.member_id = d.member_id").
		Select("d.*,um.nickname,um.avatar,um.gender,um.birthdate").
		Order("dynamic_id desc").
		Scopes(Paginate(page, prepage)).Find(&lists)
	return
}
