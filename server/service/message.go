package service

import (
	"fmt"
	"free-im/server/model"
)

// client auth handle
func ClientAuth(ctx *model.Context) {
	//认证 ctx.Message.AccessToken	&& ctx.Message.UserID


}

//client send message handle
func ClientSendMessage(ctx *model.Context) {
	fmt.Println("client send message")
}

//client pull message handle
func ClientPullMessage(ctx *model.Context) {

}


