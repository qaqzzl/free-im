package service

import (
	"database/sql"
	"fmt"
	"free-im/dao"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strconv"
	"time"
)

type AccountService struct {
}

func NewAccountService() (s *AccountService) {
	return s
}

// 账号登录
func (s *AccountService) PhoneLogin(phone string) (member_id int64, err error) {
	user_auths,_ := dao.NewMysql().Table("user_auths").
		Where(fmt.Sprintf("identifier = '%s' and identity_type = '%s'", phone, "phone")).
		First("member_id")
	member_id_int,_ := strconv.Atoi(user_auths["member_id"])
	member_id = int64(member_id_int)
	return member_id, err
}

// 账号注册
func (s *AccountService) Register(identifier string, identity_type string, credential string) (member_id int64, err error) {
	// 创建用户
	user_member := make(map[string]string)
	rand.Seed(time.Now().Unix())
	user_member["id"] = fmt.Sprintf("%06d",rand.Int31n(10000))
	user_member["nickname"] = "会员 - " + fmt.Sprintf("%06d",rand.Int31n(10000))
	user_member["avatar"] = "https://blog.cdn.qaqzz.com/icon.png"
	user_member["created_at"] = strconv.Itoa(int(time.Now().Unix()))
	user_member["updated_at"] = user_member["created_at"]
	var result sql.Result
	if result,err = dao.NewMysql().Table("user_member").Insert(user_member); err != nil {
		return member_id, err
	}
	// 创建用户账号
	user_auths := make(map[string]string)
	member_id,_ = result.LastInsertId()
	user_auths["member_id"] = strconv.Itoa(int(member_id))
	user_auths["identity_type"] = identity_type
	user_auths["identifier"] = identifier
	user_auths["credential"] = credential
	if result,err = dao.NewMysql().Table("user_auths").Insert(user_auths); err != nil {
		return member_id, err
	}

	// 返回 id
	return member_id, err
}

// 判断账号是否注册
func (s *AccountService) IsRegister(account string, types string) (bool, error) {
	c, err :=  dao.NewMysql().Table("user_auths").Where(`identity_type = '`+types+`' and identifier = '`+ account + `'`).Count()
	if c > 0 {
		return true, err
	}
	return false, err
}


// 获取用户 token
func (s *AccountService) GetToken(member_id int64, client string) (token string, err error) {
	token = uuid.NewV4().String()
	user_auths_token := make(map[string]string)
	user_auths_token["member_id"] = strconv.Itoa(int(member_id))
	user_auths_token["token"] = token
	user_auths_token["client"] = client
	user_auths_token["last_time"] = strconv.Itoa(int(time.Now().Unix()))
	user_auths_token["created_at"] = user_auths_token["last_time"]
	_,err = dao.NewMysql().Table("user_auths_token").Insert(user_auths_token)
	return token, err
}