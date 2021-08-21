package service

import (
	"errors"
	"fmt"
	"free-im/internal/http_app/dao"
	"free-im/internal/http_app/model"
	http2 "free-im/pkg/http"
	"free-im/pkg/service/id"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type AccountService struct {
}

var UserDao = dao.User

func (s *AccountService) Login(user_auths model.UserAuths, user_member model.UserMember, device_id, client string) (interface{}, error) {
	var (
		err       error
		member_id int64
		token     string
	)

	//判断是否已经注册
	result := dao.Dao.DB().Table(user_auths.TableName()).
		Where("identifier = ? and identity_type = ?", user_auths.Identifier, user_auths.IdentityType).
		Select("member_id").Find(&user_auths)
	if result.Error != nil {
		return nil, errors.New("系统忙，请稍后再试")
	}
	if user_auths.MemberId == 0 {
		member_id, err = s.Register(user_auths, user_member)
		if err != nil {
			return nil, err
		}
	} else {
		member_id = user_auths.MemberId
	}
	// 获取token
	if token, err = s.GetToken(member_id, device_id, client); err != nil {
		return nil, err
	}
	ret := make(map[string]interface{})
	ret["access_token"] = token
	ret["uid"] = member_id
	return ret, nil
}

// * 账号注册
// * identifier 账号
// * identity_type 账号类型
// * credential 密码凭证
func (s *AccountService) Register(user_auths model.UserAuths, user_member model.UserMember) (member_id int64, err error) {
	var errCreateFail = errors.New("创建账户失败")
	// 创建用户
	rand.Seed(time.Now().Unix())
	freeid, err := id.FreeID.GetID()
	if err != nil {
		return
	}
	user_member.Id = freeid
	if user_member.Nickname == "" {
		user_member.Nickname = "会员-" + freeid
	} else {
		user_member.Nickname += "-" + freeid
	}
	if user_member.Avatar == "" {
		user_member.Avatar = "http://free-im-qn.qaqzz.com/default_avatar.png"
	}
	err = dao.Dao.DB().Transaction(func(tx *gorm.DB) error {
		if dao.Dao.DB().Table(user_member.TableName()).Create(&user_member).Error != nil {
			return errCreateFail
		}
		// 创建用户账号
		user_auths.MemberId = user_member.MemberId
		if dao.Dao.DB().Table("user_auths").Create(&user_auths).Error != nil {
			return errCreateFail
		}
		return nil
	})
	if err != nil {
		return member_id, errCreateFail
	}
	member_id = user_member.MemberId
	// 添加默认好友 id=1
	timeUnix := time.Now().Unix()
	sql := fmt.Sprintf("INSERT INTO `user_friend` (member_id,friend_id,status,created_at) VALUES (%d,%d,%d,%d) "+
		"ON DUPLICATE KEY UPDATE status=VALUES(status)",
		member_id, 1, 0, timeUnix)
	dao.Dao.DB().Exec(sql)
	sql = fmt.Sprintf("INSERT INTO `user_friend` (member_id,friend_id,status,created_at) VALUES (%d,%d,%d,%d) "+
		"ON DUPLICATE KEY UPDATE status=VALUES(status)",
		1, member_id, 0, timeUnix)
	dao.Dao.DB().Exec(sql)

	// 加入默认群组
	var groupMember = model.GroupMember{
		MemberId: member_id, GroupId: 1,
		MemberIdentity: "common", Status: "normal",
		NotifyLevel: 0,
	}
	dao.Chatroom.JoinGroup(&groupMember, 1007212)

	// 返回 id
	return member_id, err
}

// 绑定推送ID
func (s *AccountService) BindPushID(member_id int64, push_id string) (err error) {
	return UserDao.BindPushID(member_id, push_id)
}

// 获取用户 token
func (s *AccountService) GetToken(member_id int64, device_id, client string) (token string, err error) {
	token, err = http2.GetToken(member_id, device_id, time.Now().Unix()+3600*12, "")
	return token, err
}
