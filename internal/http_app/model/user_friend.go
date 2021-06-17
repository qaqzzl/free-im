package model

// 用户好友表(好友申请也是这个表)
type UserFriend struct {
	Id           int64  `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	MemberId     int64  `gorm:"column:member_id;type:bigint(20);NOT NULL" json:"member_id"`                         // 会员ID
	FriendId     int64  `gorm:"column:friend_id;type:bigint(20);NOT NULL" json:"friend_id"`                         // 好友ID
	FriendRemark string `gorm:"column:friend_remark;type:varchar(50);NOT NULL" json:"friend_remark"`                // 昵称备注
	Status       int    `gorm:"column:status;type:tinyint(1);default:0;NOT NULL" json:"status"`                     // 0-正常, 1-删除
	CreatedAt    int    `gorm:"autoCreateTime;column:created_at;type:int(11);default:0;NOT NULL" json:"created_at"` // 添加时间
}

func (m *UserFriend) TableName() string {
	return "user_friend"
}
