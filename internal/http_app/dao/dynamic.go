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
