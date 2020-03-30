package service

import (
	"fmt"
	"free-im/dao"
	"strconv"
	"strings"
	"time"
)

type UserService struct {

}

// 获取会员信息
func (s *UserService) GetMemberInfo(member_id string) (user_member map[string]string, err error) {
	user_member = make(map[string]string)
	user_member,err =  dao.NewMysql().Table("user_member").
		Where("member_id = " + member_id).
		First("member_id,nickname,gender,birthdate,avatar,signature,city,province")
	return user_member,err
}

// 修改会员信息
func (s *UserService) UpdateMemberInfo(member_id string, data map[string]string) (err error) {
	var update string
	for k,v := range data {
		update +="`"+ k+"`='"+v+"',"
	}
	update = strings.Trim(update,",")
	_,err = dao.NewMysql().Table("user_member").
		Where("member_id = " + member_id).Update(update);
	return err
}

// 判断用户昵称是否存在
func (s *UserService) IsMemberNickname(member_id string, nickname string) bool {
	c, _ := dao.NewMysql().Table("user_member").
		Where(fmt.Sprintf("member_id != %s and nickname = %s",member_id, nickname) ).Count()
	if c > 0 {
		return true
	}
	return false
}

// 添加好友
func (s *UserService) AddFriend(member_id string, friend_id string, remark string) (ret map[string]string, err error) {
	ret = make(map[string]string)
	// 判断是否已经申请
	user_friend,_ := dao.NewMysql().Table("user_friend_apply").
		Where( fmt.Sprintf("member_id = %s and friend_id = %s",member_id, friend_id) ).
		Order("id desc").
		First("member_id,friend_id,status")
	if len(user_friend) == 0 {
		if err = addFriendInsert(member_id, friend_id, remark ); err != nil {
			return ret, err
		}
		ret["message"] = "已发送, 等待对方同意"
		ret["code"] = "1"
		return ret, err
	}

	switch user_friend["status"] {
	case "0":
		ret["message"] = "已发送, 等待对方同意"
		ret["code"] = "1"
		break;
	case "1":
		// 判断是否为正常好友关系
		me_is_friend,_ := queryFriendBindStatus(member_id, user_friend["friend_id"])
		to_is_friend,_ := queryFriendBindStatus(member_id, user_friend["friend_id"])
		if me_is_friend != 0 || to_is_friend != 0{
			addFriendInsert(member_id, friend_id, remark )
			ret["message"] = "已发送, 等待对方同意"
			ret["code"] = "1"
		} else {
			ret["message"] = "添加成功, 开始聊天吧"
			ret["code"] = "0"
		}
		break;
	case "2":
		if err = addFriendInsert(member_id, friend_id, remark ); err != nil {
			return ret, err
		}
		ret["message"] = "已发送, 等待对方同意"
		ret["code"] = "1"
	}
	return ret, err
}
func addFriendInsert(member_id string, friend_id string, remark string) (err error) {
	friend_insert := make(map[string]string)
	friend_insert["member_id"] = member_id
	friend_insert["friend_id"] = friend_id
	friend_insert["remark"] = remark
	friend_insert["created_at"] = strconv.Itoa(int(time.Now().Unix()))
	_, err = dao.NewMysql().Table("user_friend_apply").Insert(friend_insert)
	return err
}

// 删除好友
func (s *UserService) DelFriend(member_id string, friend_id string) (err error) {
	// 判断是否已经存在
	_,err = dao.NewMysql().Table("user_friend").
		Where( fmt.Sprintf("member_id = %s and friend_id = %s or member_id = %s and friend_id = %s",member_id, friend_id, friend_id, member_id) ).
		Delete()
	// 删除聊天室 && 聊天室消息
	return err
}

// 好友申请列表
func (s *UserService) FriendApplyList(member_id string) ([]map[string]string, error) {
	list, err := dao.NewMysql().Table("user_friend_apply as ufa").Where("ufa.friend_id = "+member_id).
		Join("INNER JOIN user_member um ON um.member_id=ufa.member_id").
		Select("um.member_id,um.nickname,um.avatar,um.signature,um.gender,ufa.remark, ufa.status, ufa.id").
		Order("id desc").
		Get()
	if len(list) == 0 {
		list = make([]map[string]string, 0)
	}
	return list, err
}

// 好友申请同意/拒绝操作
func (s *UserService) FriendApplyAction(id string, member_id string, action string) (ret map[string]string, err error) {
	ret = make(map[string]string)
	friend_apply, err := dao.NewMysql().Table("user_friend_apply").Where("id = "+id+" and status = 0").
		First("member_id,friend_id")
	if len(friend_apply) == 0 {
		ret["message"] = "操作失败"
		ret["code"] = "1"
		return ret,err
	}
	if friend_apply["friend_id"] != member_id {
		ret["message"] = "违法操作"
		ret["code"] = "1"
		return ret,err
	}
	if _,err = dao.NewMysql().Table("user_friend_apply").Where("id = "+id+" and status = 0").Update("status="+action); err != nil {
		return ret,err
	}
	if action == "1" {
		timeUnix := time.Now().Unix()
		sql := fmt.Sprintf("INSERT INTO `user_friend` (member_id,friend_id,status,created_at) VALUES (%s,%s,%d,%d) " +
			"ON DUPLICATE KEY UPDATE status=VALUES(status)",
			friend_apply["member_id"],friend_apply["friend_id"],0,timeUnix)
		if _,err = dao.MysqlConn.Exec(sql); err != nil {
			return ret,err
		}
		sql = fmt.Sprintf("INSERT INTO `user_friend` (member_id,friend_id,status,created_at) VALUES (%s,%s,%d,%d) " +
			"ON DUPLICATE KEY UPDATE status=VALUES(status)",
			friend_apply["friend_id"],friend_apply["member_id"],0,timeUnix)
		if _,err = dao.MysqlConn.Exec(sql); err != nil {
			return ret,err
		}
	}
	ret["message"] = "操作成功"
	ret["code"] = "0"
	return ret,err
}

// 好友列表
func (s *UserService) FriendList(member_id string) (list []map[string]string, err error) {
	list, err = dao.NewMysql().Table("user_friend").
		Where("member_id = "+member_id + " and status = 0").
		Select("member_id,friend_id").
		Get()
	if len(list) == 0 {
		list = make([]map[string]string, 0)
	}
	for k,vo := range list{
		var (
			where string
		)
		where = "member_id = "+ vo["friend_id"]
		member,_ := dao.NewMysql().Table("user_member").Where(where).First("avatar,gender,member_id,nickname,signature")
		list[k] = member
	}
	return list, err
}


// 搜索好友
func (s *UserService) SearchMember(search string) (list []map[string]string, err error) {
	list, err = dao.NewMysql().Table("user_member").Where("nickname = '"+search + "' or id = '"+search+"'").
		Select("member_id,nickname,avatar,signature,gender").
		Get()

	if len(list) == 0 {
		list = make([]map[string]string, 0)
	}
	return list, err
}


// 他人基本信息(他人主页)
func (s *UserService) OthersHomeInfo(member_id string, to_member_id string) (user_member map[string]string, err error) {
	user_member = make(map[string]string)
	user_member,err =  dao.NewMysql().Table("user_member").Where("member_id = " + to_member_id).First("member_id,nickname,gender,birthdate,avatar,signature,city,province")
	// 判断是否是好友关系
	user_member["is_friend"] = "no";
	is_friend,_ := queryFriendBindStatus(member_id, to_member_id)
	if is_friend == 0 {
		user_member["is_friend"] = "yes";
	}

	return user_member,err
}

// 查询好友关系状态
// return int 0:正常好友关系, 1:删除, 2: 非好友关系
func queryFriendBindStatus(member_id string, to_member_id string) (int, error) {
	user_friend, _ := dao.NewMysql().Table("user_friend").Where( fmt.Sprintf("member_id = %s and friend_id = %s",member_id, to_member_id) ).First("status")
	if len(user_friend) == 0 {
		return 2,nil
	}
	if user_friend["status"] == "1" {
		return 1,nil
	}
	return 0,nil
}