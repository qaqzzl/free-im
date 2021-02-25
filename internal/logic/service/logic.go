package service

import (
	"context"
	"encoding/json"
	"free-im/internal/logic/dao"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/rpc_client"
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
	insert := make(map[string]string)
	for _, v := range members.([]interface{}) {
		UserID := string(v.([]uint8))
		// 存储消息
		insert["message_id"] = m.MessageId
		insert["chatroom_id"] = m.ChatroomId
		insert["member_id"] = UserID
		insert["content"] = string(BodyData)
		if _, err := dao.NewMysql().Table("message").Insert(insert); err != nil {
			logger.Sugar.Error("存储消息失败", err)
			return err
		}
		//发送消息
		rpc_client.TCPConnectInit.DeliverMessageByUIDAndNotDID(ctx, &pbs.DeliverMessageReq{
			UserId:   UserID,
			DeviceId: m.DeviceID,
			Message:  &packages,
		})
	}
	// 消息回执
	packages.Action = pbs.Action_MessageACK
	packages.BodyData = []byte(m.MessageId)
	rpc_client.TCPConnectInit.DeliverMessageByUIDAndDID(ctx, &pbs.DeliverMessageReq{
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
	messages, err := dao.NewMysql().Table("message").
		Where("member_id = " + mp.UserId + " and message_id > '" + mp.MessageId + "'").
		Select("content").Get()
	if err != nil {
		logger.Sugar.Error("消息查询失败", err)
		return err
	}
	for _, v := range messages {
		packages := pbs.MsgPackage{
			Action:   pbs.Action_Message,
			BodyData: []byte(v["content"]),
		}
		//发送消息
		rpc_client.TCPConnectInit.DeliverMessageByUID(ctx, &pbs.DeliverMessageReq{
			UserId:  mp.UserId,
			Message: &packages,
		})
	}
	return nil
}
