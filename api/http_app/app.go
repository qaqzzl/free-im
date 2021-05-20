package http_app

import (
	"free-im/config"
	app_http "free-im/internal/http_app/api"
	http2 "free-im/pkg/http"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func StartHttpServer() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(config.HttpConf.HttpListenAddr)
	if err != nil {
		panic(err.Error())
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Any("", func(c *gin.Context) {
		c.String(http.StatusOK, "free-im v1.0")
	})

	// 登陆 ｜ 注册
	r.POST("/login", app_http.PhoneLogin)                        // 手机号登录 / 注册                                                            // 手机号登录 / 注册
	r.POST("/login/phone.password", app_http.PhonePasswordLogin) // 手机号密码登录 / 注册                                                            // 手机号登录 / 注册
	r.POST("/login/qq", app_http.QQLogin)                        // QQ登陆

	r.Any("/app/new.version.get", app_http.AppNewVersionGet) // http_app 最新版本获取

	authorized := r.Group("/").Use(authorizedMiddleware())
	{
		authorized.Any("/user/member.info", app_http.GetMemberInfo)                                                     // 获取会员信息
		authorized.Any("/user/update.member.info", app_http.UpdateMemberInfo)                                           // 修改会员信息
		authorized.Any("/user/others.home.info", app_http.OthersHomeInfo)                                               // 获取用户基本信息(他人主页)
		authorized.Any("/search/friend", app_http.SearchMember)                                                         // 搜索好友
		authorized.Any("/user/add.friend", app_http.AddFriend)                                                          // 添加好友
		authorized.Any("/user/del.friend", app_http.DelFriend)                                                          // 删除好友
		authorized.Any("/user/friend.apply.list", app_http.FriendApplyList)                                             // 好友申请列表
		authorized.Any("/user/friend.apply.action", app_http.FriendApplyAction)                                         // 好友申请操作
		authorized.Any("/user/friend.list", app_http.FriendList)                                                        // 好友列表
		authorized.Any("/chatroom/friend_id.get.chatroom_id", app_http.FriendIdGetChatroomId)                           // 通过好友ID 获取 聊天室ID
		authorized.Any("/chatroom/get.chatroom.avatar.name.by.chatroom_id", app_http.GetChatroomAvatarNameByChatRoomID) // 通过聊天室ID 获取 聊天室头像名称
		authorized.Any("/chatroom/chatroom.list", app_http.ChatroomList)                                                // 聊天室列表
		authorized.Any("/chatroom/create.group", app_http.CreateGroup)                                                  // 创建群组
		authorized.Any("/chatroom/add.group", app_http.AddGroup)                                                        // 加入群组
		authorized.Any("/chatroom/my.group.list", app_http.MyGroupList)                                                 // 我的群组列表
		authorized.Any("/common/get.qiniu.upload.token", app_http.GetQiniuUploadToken)                                  // 获取七牛上传token
		authorized.Any("/common/send.sms", app_http.SendLoginSms)                                                       // 发送手机号验证码
		authorized.Any("/dynamic/publish", app_http.DynamicPublish)                                                     // 发布动态
		authorized.Any("/dynamic/list", app_http.DynamicList)                                                           // 动态列表
		authorized.Any("/common/get.message.id", app_http.GetMessageId)                                                 // 获取消息ID , 临时使用

	}

	return r
}

func authorizedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			http2.Resp(c, 401, nil, "请登陆")
			c.Abort()
			return
		}
		parts := strings.SplitN(authorization, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			http2.Resp(c, 401, nil, "请求头中auth格式有误")
			c.Abort()
			return
		}
		token, err := http2.DecryptToken(parts[1])
		if err != nil {
			http2.Resp(c, 401, nil, "Token 解析失败")
			c.Abort()
			return
		}
		if token.Expire < time.Now().Unix() {
			http2.Resp(c, 401, nil, "token 已过期")
			c.Abort()
			return
		}
		c.Set("authorized_member_id", token.UserId)
		c.Set("token_info", token)

		c.Next()

	}
}
