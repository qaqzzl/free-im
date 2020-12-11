package id

import (
	"fmt"
	"free-im/pkg/protos/pbs"
	"strconv"
	"time"
)

type messageId struct {
}

var MessageID = new(messageId)

//码规则：从左至右，每 5 个 Bit 转换为一个整数，以这个整数作为下标，即可在下表中找到对应的字符。
var codingMap = [32]string{
	"2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F", "G", "H",
	"J", "K", "L", "M", "N", "P", "Q", "R",
	"S", "T", "U", "V", "W", "X", "Y", "Z",
}

// 解码
var decodeCodingMap = map[string]int{
	"2": 0, "3": 1, "4": 2, "5": 3, "6": 4, "7": 5, "8": 6, "9": 7,
	"A": 8, "B": 9, "C": 10, "D": 11, "E": 12, "F": 13, "G": 14, "H": 15,
	"J": 16, "K": 17, "L": 18, "M": 19, "N": 20, "P": 21, "Q": 22, "R": 23,
	"S": 24, "T": 25, "U": 26, "V": 27, "W": 28, "X": 29, "Y": 30, "Z": 31,
}

// 其中，自旋 ID 是一个从 0 到 4095 范围内，顺序递增的数，生成规则如下：
var max_message_seq = 0xFFF
var currentSeq = 0

func getMessageSeq() int {
	// 加锁 code...
	ret := currentSeq + 1
	if currentSeq > max_message_seq {
		currentSeq = 0
		ret = currentSeq + 1
	}
	currentSeq++
	// 解锁 code...
	return ret
}

// 获取消息ID
// chatroom_id 会话ID
// sessionType 会话类型
func (id *messageId) GetId(chatroom_id int64, sessionType pbs.ChatroomType) (message_id string) {
	// 1）获取当前系统的时间戳毫秒，并赋值给消息 ID 的高 64 Bit ：
	timeUnixNano := time.Now().UnixNano() / 1e6
	highBits := timeUnixNano

	// 2）获取一个自旋 ID ， highBits 左移 12 位，并将自旋 ID 拼接到低 12 位中：
	seq := getMessageSeq()
	highBits = highBits << 12
	highBits = highBits | int64((int64(seq))&0xFFF)
	// 3）上步的 highBits 左移 4 位，将会话类型拼接到低 4 位：
	highBits = highBits << 4
	highBits = highBits | int64((int64(sessionType))&0xF)
	// 4）取会话 ID 哈希值的低 22 位，记为 sessionIdInt：
	sessionId := chatroom_id
	sessionInt := sessionId & 0x3FFFFF

	// 5）highBits 左移 6 位，并将 sessionIdInt 的高 6 位拼接到 highBits 的低 6 位中：
	highBits = highBits << 6
	highBits = highBits | int64(sessionInt>>16)
	// 6）取会话 ID 的低 16 位作为 lowBits：
	lowBits := int64((sessionInt & 0xFFFF) << 16)
	// 7）highBits 与 lowBits 拼接得到 80 Bit 的消息 ID，对其进行 32 进制编码，即可得到唯一消息 ID：
	BitId := strconv.FormatInt(highBits, 2) + strconv.FormatInt(lowBits, 2)
	for i := 0; i < 16; i++ {
		str := BitId[i*5 : (i+1)*5]
		index, _ := strconv.ParseInt(str, 2, 0)
		message_id += codingMap[index]
	}
	return message_id
}

func (id *messageId) DecodeID(message_id string) (BitId string) {
	for i := 0; i < 16; i++ {
		index := message_id[i:(i + 1)]
		value := decodeCodingMap[index]
		strid := fmt.Sprintf("%02d", value) // 前置补0
		BitId += strid
	}
	return BitId
}

func converToBianry(n int64) string {
	result := ""
	for ; n > 0; n /= 2 {
		lsb := int(n % 2)
		result = strconv.Itoa(lsb) + result
	}
	return result
}
