package v1

import (
	"encoding/json"
	"free-im/util"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"hash/crc32"
	"net/http"
	"strconv"
	"time"
)

// 获取七牛上传token
// https://developer.qiniu.com/kodo/manual/1206/put-policy
func GetQiniuUploadToken(writer http.ResponseWriter, request *http.Request) {
	accessKey := "qW7rPngWLk8Nl3MQfehQ_G5ELAZaH47Dej50Dj7k"
	secretKey := "cN5unz025wgnfHJ_Ck3iBjpLUoByXnUVB8Uu4P1g"
	putPolicy := storage.PutPolicy{
		Scope:            "cdn1",
		//CallbackURL:      "http://api.example.com/qiniu/upload/callback",
		CallbackBody:     `{"key":"$(key)","hash":"$(etag)","mimeType":"$(mimeType)","imageInfo":$(imageInfo),"ext":"$(ext)"}`,
		CallbackBodyType: "application/json",
		FsizeLimit: 20971520,	//上传大小限制20M
		ForceSaveKey:true,		//强制使用服务端命名
		SaveKey: "free-im/test/$(etag)$(ext)",		//强制使用服务端命名
		DetectMime:1,			// 使用七牛检查 mime
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	ret := make(map[string]string)
	ret["token"] = upToken
	ret["expires"] = strconv.FormatInt(time.Now().Unix() + 3600,10)
	ret["domain"] = "https://cdn.qaqzz.com/"
	ret["message"] = "获取成功"
	ret["code"] = "0"
	util.RespOk(writer, ret, "")
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
func GetMessageId(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)

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
	sessionId := formData["chatroom_id"].(string)
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
	ret := make(map[string]interface{})
	ret["message_id"] = message_id
	util.RespOk(writer, ret, "")
}
