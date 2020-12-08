package main

import (
	"fmt"
	"free-im/internal/app/dao"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/util/id"
)

func main() {
	// 初始化 redis
	// 初始化ID 生成器
	id.Init(*dao.GetRConn())
	//value2 := converToBianry(22)
	//intvalue2,_ := strconv.Atoi(value2)
	//fmt.Println(0xFFFF)
	//fmt.Println((55 & 65535) )
	//fmt.Println((intvalue2 & 0xFFFF) << 16)
	id.ChatroomID.GetID(pbs.ChatroomType_Single)
	id.ChatroomID.GetID(pbs.ChatroomType_Single)
	id.ChatroomID.GetID(pbs.ChatroomType_Single)
	ChatroomID, _ := id.ChatroomID.GetID(pbs.ChatroomType_Single)
	fmt.Println("ChatroomID:", ChatroomID)

	memberid := id.MessageID.GetId(ChatroomID, pbs.ChatroomType_Single)
	fmt.Println(memberid)
	//fmt.Println("解码消息ID:",id.MessageID.DecodeID(memberid))

}
