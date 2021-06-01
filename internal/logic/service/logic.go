package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"free-im/internal/logic/dao"
	"free-im/internal/logic/model"
	http2 "free-im/pkg/http"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/rpc_client"
	"time"
)

// client conn auth
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

//client message handle
func MessageReceive(ctx context.Context, req pbs.MessageReceiveReq) error {
	m := req.Message
	//数据验证 code ... //

	//m.MessageSendTime = time.Now().Unix()

	BodyData, _ := json.Marshal(m)

	// 查询聊天室成员
	members, err := dao.Chatroom.GetMembers(m.ChatroomId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	// 给聊天室全员发送消息
	packages := pbs.MsgPackage{
		Action:   pbs.Action_Message,
		BodyData: BodyData,
	}
	for _, v := range members {
		UserID := v
		// 存储消息
		store_message := model.Message{
			MessageId: m.MessageId,
			Content:   string(BodyData),
		}
		store_message.ChatroomId = m.ChatroomId
		store_message.MemberId = UserID

		if err := dao.Message.StoreMessage(&store_message); err != nil {
			logger.Sugar.Error("存储消息失败", err)
			return err
		}
		//发送消息
		rpc_client.ConnectInit.DeliverMessageByUIDAndNotDID(ctx, &pbs.DeliverMessageReq{
			UserId:   UserID,
			DeviceId: m.DeviceID,
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

	return nil
}

func MessageSync(ctx context.Context, mp pbs.MessageSyncReq) error {
	var messages []model.Message
	result := dao.Dao.DB().Table("message").Where("member_id = ? and message_id > ?", mp.UserId, mp.MessageId).Select("content").Find(&messages)
	if result.Error != nil {
		logger.Sugar.Error("消息查询失败", result.Error)
		return result.Error
	}
	fmt.Println(messages)
	for _, itme := range messages {
		packages := pbs.MsgPackage{
			Action:   pbs.Action_Message,
			BodyData: []byte(itme.Content),
		}
		//发送消息
		rpc_client.ConnectInit.DeliverMessageByUID(ctx, &pbs.DeliverMessageReq{
			UserId:  mp.UserId,
			Message: &packages,
		})
	}
	return nil
}
