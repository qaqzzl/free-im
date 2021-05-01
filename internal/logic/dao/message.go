package dao

import "free-im/internal/logic/model"

type message struct {
}

var Message = new(message)

// * 储存消息
func (d *message) StoreMessage(store_message *model.Message) error {
	result := Dao.DB().Table(store_message.TableName()).Create(store_message)
	return result.Error
}
