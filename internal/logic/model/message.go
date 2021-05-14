package model

// 消息记录表
type Message struct {
	Id         int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	MessageId  string `gorm:"column:message_id;type:char(32);NOT NULL" json:"message_id"`     // 消息ID
	ChatroomId int64  `gorm:"column:chatroom_id;type:bigint(20);NOT NULL" json:"chatroom_id"` // 聊天室ID
	MemberId   int64  `gorm:"column:member_id;type:bigint(20);NOT NULL" json:"member_id"`     // 会员ID
	Content    string `gorm:"column:content;type:text" json:"content"`
}

func (m *Message) TableName() string {
	return "message"
}
