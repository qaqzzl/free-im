// Code generated by sql2gorm. DO NOT EDIT.
package model

// 群组表
type GroupMember struct {
	GroupMemberId  uint   `gorm:"column:group_member_id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"group_member_id"`
	GroupId        int    `gorm:"column:group_id;type:int(11);NOT NULL" json:"group_id"`                // 群组ID
	MemberId       int    `gorm:"column:member_id;type:int(11);NOT NULL" json:"member_id"`              // 会员ID
	MemberIdentity string `gorm:"column:member_identity;type:char(10);NOT NULL" json:"member_identity"` // 成员身份: admin-管理员, root-群主, common-普通成员
	Status         string `gorm:"column:status;type:char(10);NOT NULL" json:"status"`                   // 状态: wait-等待同意, normal-正常, refuse-拒绝, blacklist-黑名单
	CreatedAt      int    `gorm:"column:created_at;type:int(11);default:0;NOT NULL" json:"created_at"`  // 添加时间
	Group          Group  `gorm:"foreignkey:GroupId;references:GroupId"`
}

func (m *GroupMember) TableName() string {
	return "group_member"
}
