package model

// 好友申请表
type UserFriendApply struct {
	Id        int64  `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	MemberId  int64  `gorm:"column:member_id;type:bigint(20);NOT NULL" json:"member_id"`                         // 会员ID
	FriendId  int64  `gorm:"column:friend_id;type:bigint(20);NOT NULL" json:"friend_id"`                         // 好友ID
	Remark    string `gorm:"column:remark;type:varchar(50);NOT NULL" json:"remark"`                              // 添加好友备注
	Status    int    `gorm:"column:status;type:tinyint(1);default:0;NOT NULL" json:"status"`                     // 0-等待, 1-同意, 2-拒绝
	CreatedAt int    `gorm:"autoCreateTime;column:created_at;type:int(11);default:0;NOT NULL" json:"created_at"` // 添加时间
}

func (m *UserFriendApply) TableName() string {
	return "user_friend_apply"
}
