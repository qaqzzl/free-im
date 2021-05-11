package service

import (
	"free-im/internal/http_app/dao"
	"free-im/internal/http_app/model"
)

type DynamicService struct {
}

// 创建动态
func (s *DynamicService) Create(m model.Dynamic) (dynamic_id int64, err error) {
	dynamic_id, err = dao.Dynamic.Create(m)
	return
}

// 动态列表
func (s *DynamicService) DynamicList(page int, prepage int) (total int64, lists []map[string]interface{}) {
	return dao.Dynamic.DynamicList(page, prepage)
}
