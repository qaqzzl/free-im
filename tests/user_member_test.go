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

func TestIsJoinGroup(t *testing.T) {
	// 判断是否已经加入群组
	is_join, _ := dao.Dao.Ris().Do("SISMEMBER", "set_im_chatroom_member:1000222", 1)
	if is_join.(int64) == 1 {
		println("已加入", is_join.(int64))
	} else {
		println("未加入", is_join.(int64))
	}
}

func TestJoinGroup(t *testing.T) {
	res, err := dao.Dao.Ris().Do("SADD", "set_im_chatroom_member:1000222", 9999) //加入聊天室

	if err != nil {
		println(err)
	}

	if res.(int64) == 1 {
		println("成功", res.(int64))
	} else {
		println("失败", res.(int64))
	}
}
