package tcp_conn

import (
	"encoding/json"
	"free-im/dao"
	"github.com/orcaman/concurrent-map"
	"hash/crc32"
	"log"
	"net"
	"strconv"
	"time"
)

type Socket interface {
	ClientAuth(message AuthMessage)
	ClientMessage(message MessagePackage)
	Close()
}

type Logic struct {
	ctx *Context
}

// client auth handle
func (logic Logic) ClientAuth(m AuthMessage) {
	// fmt.Println("进入认证方法",m)
	//认证
	if m.UserID == "" || m.AccessToken == "" || m.DeviceID == "" || m.ClientType == "" || m.DeviceType == "" {
		return
	}
	logic.ctx.IsAuth = true
	logic.ctx.UserID = m.UserID
	logic.ctx.DeviceID = m.DeviceID
	logic.ctx.ClientType = m.ClientType
	logic.ctx.DeviceType = m.DeviceType
	clientDevice := ClientDevice{
		DeviceID:m.DeviceID,
		ClientType:m.ClientType,
		Conn: logic.ctx.TcpConn,
		Context: logic.ctx,
	}
	//加入连接集合
	if tmp, ok := SocketConnPool.Get(logic.ctx.UserID); ok {
		device_map := tmp.(cmap.ConcurrentMap)
		// 判断连接是否存在相同设备
		for k,v := range device_map.Items() {
			if k == m.DeviceType {	// 如果有同类型的设备登录了 ,通知其设备下线
				device := v.(ClientDevice)
				if device.DeviceID != m.DeviceID {
					// 通知其设备下线 code...

				}
				// 关闭连接
				device.Context.TcpConn.Close()
				for  {		// 等待上个连接完全关闭
					if device.Context.Status == false {
						break
					}
				}
			}
		}
		device_map.Set(m.DeviceType, clientDevice)
		SocketConnPool.Set(logic.ctx.UserID, device_map)
	} else  {
		device_map := cmap.New()
		device_map.Set(m.DeviceType, clientDevice)
		SocketConnPool.Set(logic.ctx.UserID, device_map)
	}

	//ctx.Response(model.Response{
	//	Code: 0,
	//	Msg: "认证成功",
	//})
}

//client send message handle
func (logic Logic) ClientMessage(m MessagePackage) {
	m.UserId = logic.ctx.UserID
	redisconn := dao.NewRedis()
	defer redisconn.Close()
	//判断是否认证 auth
	if logic.ctx.IsAuth == false {
		//ctx.ConnSocket.Close()
		return
	}
	m.MessageSendTime = time.Now().Unix()
	message,_ := json.Marshal(m)
	//字段验证 code ... //

	// <- chan
	//ctx.InChan <- &ctx.Message

	// 存储消息
	// redisconn.Do("ZADD", "sorted_set_im_chatroom_message_record:"+m.ChatroomId, message_id, message)

	// 查询聊天室成员
	members,err := redisconn.Do("SMEMBERS", "set_im_chatroom_member:"+m.ChatroomId)
	if err != nil {
		log.Println("查询聊天室成员失败", err)
		return
	}
	// 给聊天室全员发送消息
	packages := Package{
		Version: Version,
		Action: ActionMessageSend,
		BodyData: message,
	}
	for _, v := range members.([]interface {}) {
		UserID := string(v.([]uint8))

		//发送消息
		tmp, ok := SocketConnPool.Get(UserID)
		if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
			for  k, vo := range tmp.(cmap.ConcurrentMap).Items() {
				if k == logic.ctx.DeviceType && UserID == logic.ctx.UserID {		// 本设备
					continue
				}

				timeUnix := time.Now().UnixNano()  / 1e6
				device := vo.(ClientDevice)
				redisconn.Do("LPUSH", "list_message_ack_timeout_retransmit:"+UserID+":"+device.ClientType, m.MessageId)
				redisconn.Do("HSET", "hash_message_ack_timeout_retransmit:"+UserID+":"+device.ClientType, m.MessageId, strconv.FormatInt(timeUnix + 1000,10)+"|"+string(packages.BodyData))
				logic.ctx.WriteChan <- sendMessage{
					Conn:device.Conn,
					Package:packages,
				}
			}
		} else {
			// fmt.Println("设备未在线 , 未读消息写入redis")
			redisconn.Do("LPUSH", "list_message_offline:"+UserID, packages.BodyData)
		}
	}
	// 消息回执
	packages.Action = ActionMessageACK
	logic.ctx.WriteChan <- sendMessage{
		Conn:logic.ctx.TcpConn,
		Package:packages,
	}
}

func (logic Logic) ClientACK(p Package) {
	m := MessagePackage{}
	if err := json.Unmarshal(p.BodyData, &m); err != nil {
		log.Println(err.Error())
		return
	}
	redisconn := dao.NewRedis()
	defer redisconn.Close()
	redisconn.Do("HDEL", "hash_message_ack_timeout_retransmit:"+logic.ctx.UserID+":"+logic.ctx.ClientType, m.MessageId)
}


func (logic Logic) Close() {
	if logic.ctx.IsConnStatus != false {
		if user_map, ok := SocketConnPool.Get(logic.ctx.UserID); ok {
			user_map.(cmap.ConcurrentMap).Remove(logic.ctx.DeviceType)
			SocketConnPool.Set(logic.ctx.UserID,user_map)
		}
	}
	logic.ctx.IsConnStatus = false
	logic.ctx.TcpConn.Close()
}

//client pull message handle
func (logic Logic) SendResponse(conn net.Conn, p Package) (n int, err error) {
	if n, err = logic.ctx.Write(conn,p); err == nil {
		return n,nil
	} else {
		return n,err
	}
}


// 编码规则：从左至右，每 5 个 Bit 转换为一个整数，以这个整数作为下标，即可在下表中找到对应的字符。
var codingMap = [32]string {
	"2","3","4","5","6","7","8","9",
	"A","B","C","D","E","F","G","H",
	"J","K","L","M","N","P","Q","R",
	"S","T","U","V","W","X","Y","Z",
}

// 其中，自旋 ID 是一个从 0 到 4095 范围内，顺序递增的数，生成规则如下：
var max_message_seq = 0xFFF
var currentSeq = 0
func getMessageSeq() int {
	ret := currentSeq+1
	if currentSeq > max_message_seq {
		currentSeq = 0
		ret = currentSeq+1
	}
	currentSeq++
	return ret
}

// 时间戳毫秒(42),
func (logic Logic) GetMessageId(m MessagePackage) {

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
	highBits = highBits | int64(sessionType & 0xF)
	// 4）取会话 ID 哈希值的低 22 位，记为 sessionIdInt：4194304‬
	sessionId := m.ChatroomId
	sessionInt := int(crc32.ChecksumIEEE([]byte(sessionId))) & 0x3FFFFF

	// 5）highBits 左移 6 位，并将 sessionIdInt 的高 6 位拼接到 highBits 的低 6 位中：
	highBits = highBits << 6
	highBits = highBits | int64(sessionInt >> 16)
	// 6）取会话 ID 的低 16 位作为 lowBits：
	lowBits := int64((sessionInt & 0xFFFFF) << 16)
	// 7）highBits 与 lowBits 拼接得到 80 Bit 的消息 ID，对其进行 32 进制编码，即可得到唯一消息 ID：
	BitId := strconv.FormatInt(highBits, 2) + strconv.FormatInt(lowBits, 2)

	var message_id string
	for i:=0; i<16; i++ {
		str := BitId[i*5:(i+1)*5]
		index,_ := strconv.ParseInt(str,2,0)
		message_id += codingMap[index]
	}
	p := Package{
		Version:Version,
		Action:ActionMessageID,
		BodyData: []byte(message_id),
	}
	logic.ctx.Write(logic.ctx.TcpConn,p)
}




