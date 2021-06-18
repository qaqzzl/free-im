package dao

import (
	"encoding/json"
	"free-im/internal/logic/model"
	"free-im/pkg/protos/pbs"
)

type message struct {
}

var Message = new(message)

// * 储存一条消息
func (d *message) StoreMessage(members []int64, m *pbs.MsgItem) error {
	// 存储消息
	BodyData, _ := json.Marshal(m)
	message := model.Message{
		MessageId:  m.MessageId,
		Content:    string(BodyData),
		ChatroomId: m.ChatroomId,
		MemberId:   m.UserId,
	}
	result := Dao.DB().Table(new(model.Message).TableName()).Create(&message)
	if result.Error != nil {
		return result.Error
	}

	var user_messages []model.UserMessage
	for _, v := range members {
		UserID := v
		// 存储用户消息记录（关联）
		user_message := model.UserMessage{
			MessageId: m.MessageId,
		}
		user_message.MemberId = UserID
		user_messages = append(user_messages, user_message)
	}
	result = Dao.DB().Table(new(model.UserMessage).TableName()).Create(&user_messages)
	return result.Error
}
