package main

import (
	"free-im/api/http_app"
	"free-im/internal/http_app/dao"
	"free-im/pkg/util/id"
)

func main() {
	// 初始化ID 生成器
	id.Init(dao.NewRedisPool())

	http_app.StartHttpServer()
}
