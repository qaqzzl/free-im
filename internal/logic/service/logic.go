package service

import (
	"context"
	"encoding/json"
	"errors"
	"free-im/internal/logic/dao"
	"free-im/internal/logic/model"
	http2 "free-im/pkg/http"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/rpc_client"
	"free-im/pkg/service/user"
	"time"
)

// 连接认证
func TokenAuth(ctx context.Context, req pbs.TokenAuthReq) (*pbs.TokenAuthResp, error) {
	m := req.Message
	if m.UserID == 0 || m.AccessToken == "" || m.DeviceID == "" || m.ClientType == "" || m.DeviceType == "" {
		return &pbs.TokenAuthResp{Statu: false}, nil
	}
	token, err := http2.DecryptToken(m.AccessToken)
	if err != nil {
		return &pbs.TokenAuthResp{Statu: false}, errors.New("Token 解析失败")
	}
	if token.Expire < time.Now().Unix() {
		// return &pbs.TokenAuthResp{Statu: false}, errors.New("Token 已过期")
	}
	return &pbs.TokenAuthResp{Statu: true}, nil
}

// 消息处理
func MessageReceive(ctx context.Context, req pbs.MessageReceiveReq) error {
	m := req.Message
	//数据验证 code ... //

	// 查询聊天室成员
	members, err := dao.Chatroom.GetMembers(m.ChatroomId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	// 储存消息
	if dao.Message.StoreMessage(members, m) != nil {
		logger.Sugar.Error(err)
		return err
	}

	// 给聊天室全员发送消息
	m.MessageSendTime = time.Now().Unix()
	BodyData, _ := json.Marshal(m)
	packages := pbs.MsgPackage{
		Action:   pbs.Action_Message,
		BodyData: BodyData,
	}
	for _, v := range members {
		UserID := v
		// 发送消息,(排除当前设备)
		sendMessage(ctx, &pbs.DeliverMessageReq{
			UserId:   UserID,
			DeviceId: m.DeviceID,
			Message:  &packages,
		}, "uid_and_not_did")
	}

	// 消息回执
	packages.Action = pbs.Action_MessageACK
	packages.BodyData = []byte(m.MessageId)
	sendMessage(ctx, &pbs.DeliverMessageReq{
		UserId:   m.UserId,
		DeviceId: m.DeviceID,
		Message:  &packages,
	}, "uid_and_did")
	return nil
}

// 消息回执
func MessageACK(ctx context.Context, mp pbs.MessageACKReq) error {

	return nil
}

// 消息同步
func MessageSync(ctx context.Context, mp pbs.MessageSyncReq) error {
	var user_messages []model.UserMessage
	result := dao.Dao.DB().Table("user_message").Where("member_id = ? and message_id > ?", mp.UserId, mp.MessageId).Select("message_id").Find(&user_messages)
	if result.Error != nil {
		logger.Sugar.Error("消息查询失败", result.Error)
		return result.Error
	}
	var message_ids []string
	for _, itme := range user_messages {
		message_ids = append(message_ids, itme.MessageId)
	}

	var messages []model.Message
	result = dao.Dao.DB().Table("message").Where("message_id in ?", message_ids).Select("content").Find(&messages)
	if result.Error != nil {
		logger.Sugar.Error("消息查询失败", result.Error)
		return result.Error
	}
	for _, itme := range messages {
		packages := pbs.MsgPackage{
			Action:   pbs.Action_Message,
			BodyData: []byte(itme.Content),
		}
		//发送消息
		sendMessage(ctx, &pbs.DeliverMessageReq{
			UserId:  mp.UserId,
			Message: &packages,
		}, "uid")
	}
	return nil
}

func sendMessage(ctx context.Context, msg *pbs.DeliverMessageReq, bykey string) (*pbs.DeliverMessageResp, error) {
	if s, err := user.User.GetUserOnline(msg.UserId); err == nil && s {
		//发送消息
		switch bykey {
		case "uid":
			return rpc_client.ConnectInit.DeliverMessageByUID(ctx, msg)
		case "uid_and_did":
			return rpc_client.ConnectInit.DeliverMessageByUIDAndDID(ctx, msg)
		case "uid_and_not_did":
			return rpc_client.ConnectInit.DeliverMessageByUIDAndNotDID(ctx, msg)
		}
	}
	return nil, nil
}
