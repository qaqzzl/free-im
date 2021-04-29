package service

import (
	"database/sql"
	"errors"
	"fmt"
	"free-im/internal/http_app/dao"
	"free-im/internal/http_app/model"
	"free-im/pkg/util/id"
	"github.com/satori/go.uuid"
	"math/rand"
	"strconv"
	"time"
)

type AccountService struct {
}

func (s *AccountService) Login(identifier string, identity_type string, credential string, data map[string]string) (interface{}, error) {
	//判断是否已经注册
	var user_auths model.UserAuths
	result := dao.Dao.DB().Table(user_auths.TableName()).
		Where(fmt.Sprintf("identifier = '%s' and identity_type = '%s'", identifier, identity_type)).
		Select("member_id").First(&user_auths)
	if result.Error != nil {
		return nil, errors.New("系统忙，请稍后再试")
	}
	var member_id uint
	if user_auths.MemberId != 0 {
		member_id = user_auths.MemberId
	} else {
		member_id, err := s.Register(identifier, identity_type, credential, data)
		if err != nil {
			return nil, err
		}
	}
	// 获取token
	var token string
	if token, err := s.GetToken(member_id, "android"); err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	ret["access_token"] = token
	ret["uid"] = member_id
	return ret, nil
}

// 账号注册
func (s *AccountService) Register(identifier string, identity_type string, credential string, user_member map[string]string) (member_id uint, err error) {
	// 创建用户
	if user_member == nil {
		user_member = make(map[string]string)
	}
	rand.Seed(time.Now().Unix())
	freeid, err := id.FreeID.GetID()
	if err != nil {
		return
	}
	user_member["id"] = freeid
	if user_member["nickname"] == "" {
		user_member["nickname"] = "会员-" + freeid
	} else {
		user_member["nickname"] += "-" + freeid
	}
	if user_member["avatar"] == "" {
		user_member["avatar"] = "http://free-im-qn.qaqzz.com/default_avatar.png"
	}
	user_member["created_at"] = strconv.Itoa(int(time.Now().Unix()))
	user_member["updated_at"] = user_member["created_at"]
	var result sql.Result
	if result, err = dao.NewMysql().Table("user_member").Insert(user_member); err != nil {
		return member_id, err
	}
	// 创建用户账号
	user_auths := make(map[string]string)
	LastInsertId, _ := result.LastInsertId()
	member_id = uint(LastInsertId)
	user_auths["member_id"] = member_id
	user_auths["identity_type"] = identity_type
	user_auths["identifier"] = identifier
	user_auths["credential"] = credential
	if result, err = dao.NewMysql().Table("user_auths").Insert(user_auths); err != nil {
		return member_id, err
	}

	// 自动添加好友 id=1
	timeUnix := time.Now().Unix()
	sql := fmt.Sprintf("INSERT INTO `user_friend` (member_id,friend_id,status,created_at) VALUES (%s,%s,%d,%d) "+
		"ON DUPLICATE KEY UPDATE status=VALUES(status)",
		member_id, "1", 0, timeUnix)
	dao.MysqlConn.Exec(sql)
	sql = fmt.Sprintf("INSERT INTO `user_friend` (member_id,friend_id,status,created_at) VALUES (%s,%s,%d,%d) "+
		"ON DUPLICATE KEY UPDATE status=VALUES(status)",
		"1", member_id, 0, timeUnix)
	dao.MysqlConn.Exec(sql)

	// 返回 id
	return member_id, err
}

// 获取用户 token
func (s *AccountService) GetToken(member_id uint, client string) (token string, err error) {
	token = uuid.NewV4().String()
	user_auths_token := make(map[string]string)
	user_auths_token["member_id"] = member_id
	user_auths_token["token"] = token
	user_auths_token["client"] = client
	user_auths_token["last_time"] = strconv.Itoa(int(time.Now().Unix()))
	user_auths_token["created_at"] = user_auths_token["last_time"]
	_, err = dao.NewMysql().Table("user_auths_token").Insert(user_auths_token)
	return token, err
}
