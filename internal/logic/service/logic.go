package service

import (
	"context"
	"encoding/json"
	"free-im/internal/logic/dao"
	"free-im/internal/logic/model"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/rpc_client"
	"strconv"
)

// client conn auth
func TokenAuth(ctx context.Context, req pbs.TokenAuthReq) (*pbs.TokenAuthResp, error) {
	m := req.Message
	if m.UserID == "" || m.AccessToken == "" || m.DeviceID == "" || m.ClientType == "" || m.DeviceType == "" {
		return &pbs.TokenAuthResp{Statu: false}, nil
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
	members, err := dao.GetRConn().Do("SMEMBERS", "set_im_chatroom_member:"+m.ChatroomId)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	// 给聊天室全员发送消息
	packages := pbs.MsgPackage{
		Action:   pbs.Action_Message,
		BodyData: BodyData,
	}
	for _, v := range members.([]interface{}) {
		UserID := string(v.([]uint8))
		// 存储消息
		store_message := model.Message{
			MessageId: m.MessageId,
			Content:   string(BodyData),
		}
		store_message.ChatroomId, _ = strconv.Atoi(m.ChatroomId)
		store_message.MemberId, _ = strconv.Atoi(UserID)

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
	var messages []*model.Message
	result := dao.Dao.DB().Table("message").Where("member_id = ? and message_id > ?", mp.UserId, mp.MessageId).Select("content").Find(&messages)
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
		rpc_client.ConnectInit.DeliverMessageByUID(ctx, &pbs.DeliverMessageReq{
			UserId:  mp.UserId,
			Message: &packages,
		})
	}
	return nil
}
