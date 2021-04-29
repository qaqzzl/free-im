package dao

import (
	"fmt"
	"free-im/internal/http_app/model"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

type chatroom struct {
}

var Chatroom = new(chatroom)

func (d *chatroom) GetGroupByID(id string, selects ...string) (group model.Group, err error) {
	Dao.DB().Table(group.TableName()).Select(selects).Where("id = ?", "id").Find(&group)
	return
}

// * 获取聊天室成员
func (d *chatroom) GetMembers(chatroom_id string) (member_ids []uint, err error) {
	members, err := Dao.Ris().Do("SMEMBERS", "set_im_chatroom_member:"+chatroom_id)
	for _, v := range members.([]interface{}) {
		user_id, _ := strconv.Atoi(string(v.([]uint8)))
		member_ids = append(member_ids, uint(user_id))
	}
	return
}

// * 创建群组
func (d *chatroom) CreateGroup(group model.Group) (group_id uint, err error) {
	rand.Seed(time.Now().Unix())
	group.Id = fmt.Sprintf("%06d", rand.Int31n(10000))
	result := Dao.DB().Table("`group`").Create(&group)
	group_id = group.GroupId
	if err = result.Error; err != nil {
		return
	}
	// 群组成员
	group_member := model.GroupMember{
		GroupId:        group_id,
		MemberId:       group.FounderMemberId,
		MemberIdentity: "root",
		Status:         "normal",
	}
	Dao.DB().Table("`group_member`").Create(&group_member)
	// redis
	_, err = Dao.Ris().Do("SADD", "set_im_chatroom_member:"+group.ChatroomId, group.OwnerMemberId) //创建聊天室
	return
}

// * 加入群组
func (d *chatroom) JoinGroup(chatroom_id string, member_id uint) (err error) {
	_, err = Dao.Ris().Do("SADD", "set_im_chatroom_member:"+chatroom_id, member_id) //加入聊天室
	return
}

// * 群组是否存在
func (d *chatroom) GroupIsExistByID(id string) (is bool, err error) {
	var c int64
	result := Dao.DB().Table("`group`").Where("id = ?", "id").Count(&c)
	err = result.Error
	if c > 0 {
		is = true
	} else {
		is = false
	}
	return
}

// * 会员群组列表 by member_id
func (d *chatroom) MemberGroupListByUID(member_id uint) (MemberGroups []*model.GroupMember, err error) {
	return d.MemberGroupList("member_id = ?", member_id)
}

// * 会员群组列表
func (d *chatroom) MemberGroupList(query interface{}, args ...interface{}) (MemberGroups []*model.GroupMember, err error) {
	model := model.GroupMember{}
	result := Dao.DB().Table(model.TableName()).Preload("Group", func(db *gorm.DB) *gorm.DB {
		return db.Select("*")
	}).Where(query, args...).Find(&MemberGroups)
	err = result.Error
	return
}