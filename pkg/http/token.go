package http

import (
	"encoding/base64"
	"encoding/json"
	"free-im/pkg/logger"
	"free-im/pkg/util"
	jsoniter "github.com/json-iterator/go"
)

type TokenInfo struct {
	UserId   int64  `json:"user_id"`   // 用户id
	DeviceId string `json:"device_id"` // 设备id
	Extend   string `json:"extend"`    // 扩展
	Expire   int64  `json:"expire"`    // 过期时间
}

// GetToken 获取token
func GetToken(userId int64, deviceId string, expire int64, extend string) (string, error) {
	info := TokenInfo{
		UserId:   userId,
		DeviceId: deviceId,
		Extend:   extend,
		Expire:   expire,
	}
	bytes, err := json.Marshal(info)
	if err != nil {
		logger.Sugar.Error(err)
		return "", err
	}

	token, err := util.RsaEncrypt(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}

// DecryptToken 对加密的token进行解码
func DecryptToken(token string) (*TokenInfo, error) {
	bytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	result, err := util.RsaDecrypt(bytes)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	var info TokenInfo
	err = jsoniter.Unmarshal(result, &info)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return &info, nil
}

// RefreshToken
