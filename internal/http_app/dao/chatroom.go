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

func (d *chatroom) GetGroupByID(group_id int64, selects ...string) (group model.Group, err error) {
	Dao.DB().Table(group.TableName()).Select(selects).Where("group_id = ?", group_id).Find(&group)
	return
}

// * 获取聊天室成员
func (d *chatroom) GetMembers(chatroom_id int64) (member_ids []int64, err error) {
	members, err := Dao.Ris().Do("SMEMBERS", "set_im_chatroom_member:"+strconv.Itoa(int(chatroom_id)))
	for _, v := range members.([]interface{}) {
		user_id, _ := strconv.Atoi(string(v.([]uint8)))
		member_ids = append(member_ids, int64(user_id))
	}
	return
}

// * 创建群组
func (d *chatroom) CreateGroup(group *model.Group) (err error) {
	rand.Seed(time.Now().Unix())
	group.Id = fmt.Sprintf("%06d", (rand.Int31n((999999 - 10000)) - 10000))
	result := Dao.DB().Table("`group`").Create(&group)
	if err = result.Error; err != nil {
		return
	}
	// 群组成员
	group_member := model.GroupMember{
		GroupId:        group.GroupId,
		MemberId:       group.FounderMemberId,
		MemberIdentity: "root",
		Status:         "normal",
		NotifyLevel:    0,
	}
	Dao.DB().Table("`group_member`").Create(&group_member)
	// redis
	_, err = Dao.Ris().Do("SADD", "set_im_chatroom_member:"+strconv.Itoa(int(group.ChatroomId)), group.OwnerMemberId) //创建聊天室
	return err
}

// * 加入群组
func (d *chatroom) JoinGroup(m *model.GroupMember, chatroom_id int64) (err error) {
	_, err = Dao.Ris().Do("SADD", "set_im_chatroom_member:"+strconv.Itoa(int(chatroom_id)), m.MemberId) //加入聊天室
	if err == nil {
		reslut := Dao.DB().Create(&m)
		err = reslut.Error
	}
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
func (d *chatroom) MemberGroupListByUID(member_id int64) (Groups []model.Group, err error) {
	var groupMember model.GroupMember
	var MemberGroups []*model.GroupMember
	Dao.DB().Table(groupMember.TableName()).Where("member_id = ?", member_id).Find(&MemberGroups)
	var group_ids []int64
	for _, v := range MemberGroups {
		group_ids = append(group_ids, v.GroupId)
	}
	if len(group_ids) > 0 {
		var group model.Group
		result := Dao.DB().Table(group.TableName()).Where("group_id in ?", group_ids).Find(&Groups)
		err = result.Error
	}
	return
}

// * 群组信息
func (d *chatroom) GroupInfo(group_id int64) (group model.Group, err error) {
	result := Dao.DB().Table(group.TableName()).Where("group_id = ?", group_id).Find(&group)
	err = result.Error
	return
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

var d = new(model.GroupMember)
