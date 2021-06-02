package http_app

import (
	"free-im/config"
	app_http "free-im/internal/http_app/api"
	http2 "free-im/pkg/http"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
	r.POST("/login", app_http.PhoneLogin)             // 手机号登录 / 注册                                                            // 手机号登录 / 注册
	r.POST("/login/qq", app_http.QQLogin)             // QQ登陆
	r.POST("/common/send.sms", app_http.SendLoginSms) // 发送手机号验证码

	r.Any("/app/new.version.get", app_http.AppNewVersionGet) // http_app 最新版本获取

	authorized := r.Group("/").Use(authorizedMiddleware())
	{
		authorized.POST("/account/bind.push_id/:push_id", app_http.BindPushID)                 // 绑定推送ID
		authorized.POST("/user/member.info", app_http.GetMemberInfo)                           // 获取会员信息
		authorized.POST("/user/update.member.info", app_http.UpdateMemberInfo)                 // 修改会员信息
		authorized.POST("/user/others.home.info", app_http.OthersHomeInfo)                     // 获取用户基本信息(他人主页)
		authorized.POST("/search/friend", app_http.SearchMember)                               // 搜索好友
		authorized.POST("/search/group", app_http.SearchGroup)                                 // 搜索好友
		authorized.POST("/user/add.friend", app_http.AddFriend)                                // 添加好友
		authorized.POST("/user/del.friend", app_http.DelFriend)                                // 删除好友
		authorized.POST("/user/friend.apply.list", app_http.FriendApplyList)                   // 好友申请列表
		authorized.POST("/user/friend.apply.action", app_http.FriendApplyAction)               // 好友申请操作
		authorized.POST("/user/friend.list", app_http.FriendList)                              // 好友列表
		authorized.POST("/chatroom/friend_id.get.chatroom_id", app_http.FriendIdGetChatroomId) // 通过好友ID 获取 聊天室ID
		authorized.POST("/chatroom/get.chatroom.info", app_http.GetChatroomBaseInfo)           // 通过聊天室ID获取聊天室信息
		authorized.POST("/chatroom/chatroom.list", app_http.ChatroomList)                      // 聊天室列表
		authorized.POST("/chatroom/create.group", app_http.CreateGroup)                        // 创建群组
		authorized.POST("/chatroom/add.group", app_http.AddGroup)                              // 加入群组
		authorized.GET("/chatroom/my.group.list", app_http.MyGroupList)                        // 我的群组列表
		authorized.POST("/chatroom/group.info", app_http.GroupInfo)                            // 群组信息
		authorized.GET("/chatroom/group.member/:group_id", app_http.GroupMember)               // 我的群组列表
		authorized.POST("/chatroom/add.group.member", app_http.AddGroupMember)                 // 添加群组成员
		authorized.POST("/common/get.qiniu.upload.token", app_http.GetQiniuUploadToken)        // 获取七牛上传token
		authorized.POST("/dynamic/publish", app_http.DynamicPublish)                           // 发布动态
		authorized.POST("/dynamic/list", app_http.DynamicList)                                 // 动态列表
		authorized.POST("/common/get.message.id", app_http.GetMessageId)                       // 获取消息ID
		authorized.POST("/message/push.message", app_http.PushMessage)                         // 发送消息
	}

	return r
}

func authorizedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			http2.Resp(c, http2.HTTP_CODE_ACCOUNT_TOKEN_ERROR, nil, "请登陆")
			c.Abort()
			return
		}
		parts := strings.SplitN(authorization, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			http2.Resp(c, http2.HTTP_CODE_ACCOUNT_TOKEN_ERROR, nil, "请求头中auth格式有误")
			c.Abort()
			return
		}
		token, err := http2.DecryptToken(parts[1])
		if err != nil {
			http2.Resp(c, http2.HTTP_CODE_ACCOUNT_TOKEN_ERROR, nil, "Token 解析失败")
			c.Abort()
			return
		}
		//if token.Expire < time.Now().Unix() {
		//	http2.Resp(c, http2.HTTP_CODE_ACCOUNT_TOKEN_ERROR, nil, "token 已过期")
		//	c.Abort()
		//	return
		//}
		c.Set("authorized_member_id", token.UserId)
		c.Set("token_info", token)

		c.Next()

	}
}
