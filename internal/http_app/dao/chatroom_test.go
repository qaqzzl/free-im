package dao

import (
	"fmt"
	"free-im/internal/http_app/model"
	"math/rand"
	"testing"
)

func TestCreateGroup(t *testing.T) {
	group := model.Group{
		Name:            "test",
		Avatar:          "http://avatar.com",
		ChatroomId:      "1",
		OwnerMemberId:   "1",
		FounderMemberId: "1",
	}
	group.Id = fmt.Sprintf("%06d", rand.Int31n(10000))
	fmt.Println(group)
	id, err := Chatroom.CreateGroup(group)
	fmt.Println(id, err)
}
