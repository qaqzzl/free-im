package tcp_conn

import (
	"context"
	"fmt"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/rpc_client"
	"testing"
)

func TestRpcLogicInit(t *testing.T) {
	// 初始化 rpc 客户端
	rpc_client.InitLogic("127.0.0.1:5000")

	m := pbs.MsgAuth{
		DeviceID:    "1",
		UserID:      2,
		AccessToken: "3",
		DeviceType:  "mobile",
		ClientType:  "android",
	}
	if resp, err := rpc_client.LogicInit.TokenAuth(context.TODO(), &pbs.TokenAuthReq{
		Message: &m,
	}); err != nil {
		logger.Sugar.Error(err)
		return
	} else if resp.Statu == false {
		return
	} else {
		fmt.Println(resp)
	}
}
