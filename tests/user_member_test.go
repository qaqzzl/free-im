package main

import (
	"free-im/internal/http_app/dao"
	"free-im/internal/http_app/model"
	"testing"
)

func TestCreatedAt(t *testing.T) {
	var user_member model.UserMember

	if user_member.Nickname == "" {
		user_member.Nickname = "会员-" + "test"
	} else {
		user_member.Nickname += "-" + "test"
	}
	if user_member.Avatar == "" {
		user_member.Avatar = "http://free-im-qn.qaqzz.com/default_avatar.png"
	}
	dao.Dao.DB().Table(user_member.TableName()).Create(&user_member)
}
