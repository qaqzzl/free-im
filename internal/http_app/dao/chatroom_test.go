package dao

import (
	"fmt"
	"testing"
)

func TestChatroom_CreateGroup(t *testing.T) {
	//group := model.Group{
	//	Name:            "test",
	//	Avatar:          "http://avatar.com",
	//	ChatroomId:      1,
	//	OwnerMemberId:   1,
	//	FounderMemberId: 1,
	//}
	//group.Id = fmt.Sprintf("%06d", rand.Int31n(10000))
	//fmt.Println(group)
	//id, err := Chatroom.CreateGroup(group)
	//fmt.Println(id, err)
}

func TestChatroom_GetGroupByID(T *testing.T) {
	group, err := Chatroom.GetGroupByID(1, "chatroom_id")
	fmt.Println(group.ChatroomId)
	fmt.Println(err)
}

func TestChatroom_MemberGroupList(t *testing.T) {
	MemberGroups, _ := Chatroom.MemberGroupListByUID(1)
	for k, v := range MemberGroups {
		fmt.Println(k, v)
	}
}
