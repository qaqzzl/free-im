package service

import (
	"context"
	"free-im/internal/app/dao"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/rpc_client"
	"log"
)

type service struct {
}

var Service = new(service)

func (s *service) OfflineMessage(ctx context.Context, UserID string) {
	// 离线消息
	go func() {
		redisconn := dao.RedisConn()
		for {
			message, err := redisconn.Do("RPOP", "list_message_offline:"+UserID)
			if err != nil {
				log.Println("list_message_offline err", err.Error())
				break
			}
			if message == nil {
				continue
			}
			m := pbs.MessagePackage{
				Action:   pbs.Action_Message,
				BodyData: message.([]byte),
			}
			rpc_client.ConnectInit.DeliverMessageByUID(ctx, &pbs.DeliverMessageReq{
				UserId:  UserID,
				Message: &m,
			})
		}
	}()
}
