package service

import (
	"errors"
	"fmt"
	"free-im/internal/http_app/dao"
	"free-im/internal/http_app/model"
	"free-im/pkg/id"
	"github.com/satori/go.uuid"
	"math/rand"
	"time"
)

type AccountService struct {
}

func (s *AccountService) Login(identifier string, identity_type string, credential string, user_member model.UserMember) (interface{}, error) {
	var (
		err       error
		member_id uint
		token     string
	)

	//判断是否已经注册
	var user_auths model.UserAuths
	result := dao.Dao.DB().Table(user_auths.TableName()).
		Where(fmt.Sprintf("identifier = '%s' and identity_type = '%s'", identifier, identity_type)).
		Select("member_id").First(&user_auths)
	if result.Error != nil {
		return nil, errors.New("系统忙，请稍后再试")
	}
	if user_auths.MemberId != 0 {
		member_id = user_auths.MemberId
	} else {
		member_id, err = s.Register(identifier, identity_type, credential, user_member)
		if err != nil {
			return nil, err
		}
	}
	// 获取token
	if token, err = s.GetToken(member_id, "android"); err != nil {
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
func (s *AccountService) Register(identifier string, identity_type string, credential string, user_member model.UserMember) (member_id uint, err error) {
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
	result := dao.Dao.DB().Table(user_member.TableName()).Create(&user_member)
	if result.Error != nil {
		return member_id, err
	}
	// 创建用户账号
	var user_auths model.UserAuths
	user_auths.MemberId = user_member.MemberId
	user_auths.IdentityType = identity_type
	user_auths.Identifier = identifier
	user_auths.Credential = credential
	if result = dao.Dao.DB().Table("user_auths").Create(&user_auths); err != nil {
		return member_id, err
	}

	// 自动添加好友 id=1
	timeUnix := time.Now().Unix()
	sql := fmt.Sprintf("INSERT INTO `user_friend` (member_id,friend_id,status,created_at) VALUES (%s,%s,%d,%d) "+
		"ON DUPLICATE KEY UPDATE status=VALUES(status)",
		member_id, "1", 0, timeUnix)
	dao.Dao.DB().Exec(sql)
	sql = fmt.Sprintf("INSERT INTO `user_friend` (member_id,friend_id,status,created_at) VALUES (%s,%s,%d,%d) "+
		"ON DUPLICATE KEY UPDATE status=VALUES(status)",
		"1", member_id, 0, timeUnix)
	dao.Dao.DB().Exec(sql)

	// 返回 id
	return member_id, err
}

// 获取用户 token
func (s *AccountService) GetToken(member_id uint, client string) (token string, err error) {
	token = uuid.NewV4().String()
	user_auths_token := model.UserAuthsToken{
		MemberId: member_id,
		Token:    token,
		Client:   client,
		LastTime: int(time.Now().Unix()),
	}
	result := dao.Dao.DB().Table(user_auths_token.TableName()).Create(&user_auths_token)
	return token, result.Error
}
