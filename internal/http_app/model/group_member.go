// Code generated by sql2gorm. DO NOT EDIT.
package model

// 群组表
type GroupMember struct {
	GroupMemberId  int64  `gorm:"column:group_member_id;primary_key;AUTO_INCREMENT" json:"group_member_id"`
	GroupId        int64  `gorm:"column:group_id;NOT NULL" json:"group_id"`               // 群组ID
	MemberId       int64  `gorm:"column:member_id;NOT NULL" json:"member_id"`             // 会员ID
	Alias          string `gorm:"column:alias;NOT NULL" json:"alias"`                     // 会员群别名
	NotifyLevel    int    `gorm:"column:notify_level;NOT NULL" json:"notify_level"`       // 通知级别，0：正常，1：接收消息但不提醒，2：屏蔽群消息
	MemberIdentity string `gorm:"column:member_identity;NOT NULL" json:"member_identity"` // 成员身份: admin-管理员, root-群主, common-普通成员
	Status         string `gorm:"column:status;NOT NULL" json:"status"`                   // 状态: normal-正常, blacklist-黑名单
	CreatedAt      int    `gorm:"column:created_at;default:0;NOT NULL" json:"created_at"` // 添加时间
}

func (m *GroupMember) TableName() string {
	return "group_member"
}
