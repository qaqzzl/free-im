package http_app

import (
	"free-im/config"
	app_http "free-im/internal/http_app/controller"
	"net/http"
)

func StartHttpServer() {
	http.HandleFunc("/login", app_http.PhoneLogin)                                                                   // 手机号登录 / 注册
	http.HandleFunc("/login/qq", app_http.QQLogin)                                                                   // 手机号登录 / 注册
	http.HandleFunc("/user/member.info", app_http.GetMemberInfo)                                                     // 获取会员信息
	http.HandleFunc("/user/update.member.info", app_http.UpdateMemberInfo)                                           // 修改会员信息
	http.HandleFunc("/user/others.home.info", app_http.OthersHomeInfo)                                               // 获取用户基本信息(他人主页)
	http.HandleFunc("/search/friend", app_http.SearchMember)                                                         // 搜索好友
	http.HandleFunc("/user/add.friend", app_http.AddFriend)                                                          // 添加好友
	http.HandleFunc("/user/del.friend", app_http.DelFriend)                                                          // 删除好友
	http.HandleFunc("/user/friend.apply.list", app_http.FriendApplyList)                                             // 好友申请列表
	http.HandleFunc("/user/friend.apply.action", app_http.FriendApplyAction)                                         // 好友申请操作
	http.HandleFunc("/user/friend.list", app_http.FriendList)                                                        // 好友列表
	http.HandleFunc("/chatroom/friend_id.get.chatroom_id", app_http.FriendIdGetChatroomId)                           // 通过好友ID 获取 聊天室ID
	http.HandleFunc("/chatroom/get.chatroom.avatar.name.by.chatroom_id", app_http.GetChatroomAvatarNameByChatRoomID) // 通过聊天室ID 获取 聊天室头像名称
	http.HandleFunc("/chatroom/chatroom.list", app_http.ChatroomList)                                                // 聊天室列表
	http.HandleFunc("/chatroom/create.group", app_http.CreateGroup)                                                  // 创建群组
	http.HandleFunc("/chatroom/add.group", app_http.AddGroup)                                                        // 加入群组
	http.HandleFunc("/chatroom/my.group.list", app_http.MyGroupList)                                                 // 我的群组列表
	http.HandleFunc("/common/get.qiniu.upload.token", app_http.GetQiniuUploadToken)                                  // 获取七牛上传token
	http.HandleFunc("/common/send.sms", app_http.SendLoginSms)                                                       // 发送手机号验证码
	http.HandleFunc("/dynamic/publish", app_http.DynamicPublish)                                                     // 发布动态
	http.HandleFunc("/dynamic/list", app_http.DynamicList)                                                           // 动态列表
	http.HandleFunc("/http_app/new.version.get", app_http.AppNewVersionGet)                                          // http_app 最新版本获取

	http.HandleFunc("/common/get.message.id", app_http.GetMessageId) // 获取消息ID , 临时使用
	err := http.ListenAndServe(config.HttpConf.HttpListenAddr, nil)
	if err != nil {
		panic(err.Error())
	}
}
