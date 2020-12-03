package service

import (
	"context"
	"encoding/json"
	"fmt"
	"free-im/internal/app/dao"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/rpc_client"
	"hash/crc32"
	"strconv"
	"time"
)

//client message handle
func MessageReceive(ctx context.Context, req pbs.MessageReceiveReq) error {
	m := req.Message
	//数据验证 code ... //

	redisconn := dao.NewRedis()
	defer redisconn.Close()

	m.MessageSendTime = time.Now().Unix()
	BodyData, _ := json.Marshal(m)

	// 存储消息
	//_, err := redisconn.Do("ZADD", "sorted_set_im_chatroom_message_record:"+m.ChatroomId, m.MessageId, BodyData)
	//if err != nil {
	//	logger.Sugar.Error("存储消息失败",err)
	//	return err
	//}

	// 查询聊天室成员
	members, err := redisconn.Do("SMEMBERS", "set_im_chatroom_member:"+m.ChatroomId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	// 给聊天室全员发送消息
	packages := pbs.MessagePackage{
		Action:   pbs.Action_Message,
		BodyData: BodyData,
	}
	for _, v := range members.([]interface{}) {
		UserID := string(v.([]uint8))
		fmt.Println(UserID)
		//发送消息
		rpc_client.ConnectInit.DeliverMessageByUID(ctx, &pbs.DeliverMessageReq{
			UserId:   UserID,
			Message:  &packages,
		})
	}
	// 消息回执
	packages.Action = pbs.Action_MessageACK
	packages.BodyData = []byte(m.MessageId)
	rpc_client.ConnectInit.DeliverMessageByUIDAndDID(ctx, &pbs.DeliverMessageReq{
		UserId:   m.UserId,
		DeviceId: m.DeviceID,
		Message:  &packages,
	})
	return nil
}

func MessageACK(ctx context.Context, mp pbs.MessageACKReq) error {
	redisconn := dao.NewRedis()
	defer redisconn.Close()
	redisconn.Do("HDEL", "hash_message_ack_timeout_retransmit", mp.MessageId)
	return nil
}

// 编码规则：从左至右，每 5 个 Bit 转换为一个整数，以这个整数作为下标，即可在下表中找到对应的字符。
var codingMap = [32]string{
	"2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F", "G", "H",
	"J", "K", "L", "M", "N", "P", "Q", "R",
	"S", "T", "U", "V", "W", "X", "Y", "Z",
}

// 其中，自旋 ID 是一个从 0 到 4095 范围内，顺序递增的数，生成规则如下：
var max_message_seq = 0xFFF
var currentSeq = 0

func getMessageSeq() int {
	ret := currentSeq + 1
	if currentSeq > max_message_seq {
		currentSeq = 0
		ret = currentSeq + 1
	}
	currentSeq++
	return ret
}

// 时间戳毫秒(42),
func GetMessageId(ctx context.Context, mp pbs.MessagePackage) {
	m := pbs.MessageItem{}
	if err := json.Unmarshal(mp.BodyData, &m); err != nil {
		logger.Sugar.Error(err)
		return
	}
	// 1）获取当前系统的时间戳毫秒，并赋值给消息 ID 的高 64 Bit ：
	highBits := time.Now().UnixNano() / 1e6
	//highBits = 1589403510000
	// 2）获取一个自旋 ID ， highBits 左移 12 位，并将自旋 ID 拼接到低 12 位中：
	seq := getMessageSeq()
	highBits = highBits << 12
	highBits = highBits | int64(seq)
	// 3）上步的 highBits 左移 4 位，将会话类型拼接到低 4 位：
	sessionType := 1
	highBits = highBits << 4
	highBits = highBits | int64(sessionType&0xF)
	// 4）取会话 ID 哈希值的低 22 位，记为 sessionIdInt：4194304‬
	sessionId := m.ChatroomId
	sessionInt := int(crc32.ChecksumIEEE([]byte(sessionId))) & 0x3FFFFF

	// 5）highBits 左移 6 位，并将 sessionIdInt 的高 6 位拼接到 highBits 的低 6 位中：
	highBits = highBits << 6
	highBits = highBits | int64(sessionInt>>16)
	// 6）取会话 ID 的低 16 位作为 lowBits：
	lowBits := int64((sessionInt & 0xFFFFF) << 16)
	// 7）highBits 与 lowBits 拼接得到 80 Bit 的消息 ID，对其进行 32 进制编码，即可得到唯一消息 ID：
	BitId := strconv.FormatInt(highBits, 2) + strconv.FormatInt(lowBits, 2)

	var message_id string
	for i := 0; i < 16; i++ {
		str := BitId[i*5 : (i+1)*5]
		index, _ := strconv.ParseInt(str, 2, 0)
		message_id += codingMap[index]
	}
	rpc_client.ConnectInit.DeliverMessageByUIDAndDID(ctx, &pbs.DeliverMessageReq{
		UserId:   m.UserId,
		DeviceId: m.DeviceID,
		Message: &pbs.MessagePackage{
			Version:  mp.Version,
			Action:   pbs.Action_GetMessageID,
			BodyData: []byte(message_id),
		},
	})
}
