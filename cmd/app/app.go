package main

import (
	"free-im/api/app"
	"free-im/internal/app/dao"
	"free-im/pkg/util/id"
)

func main() {
	// 初始化ID 生成器
	id.Init(dao.RedisConn())

	app.StartHttpServer()
}
