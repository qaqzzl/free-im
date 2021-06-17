package model

// 用户消息记录表
type UserMessage struct {
	Id        int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	MessageId string `gorm:"column:message_id;type:char(32);NOT NULL" json:"message_id"` // 消息ID
	MemberId  int64  `gorm:"column:member_id;type:bigint(20);NOT NULL" json:"member_id"` // 会员ID
}

func (m *UserMessage) TableName() string {
	return "user_message"
}
