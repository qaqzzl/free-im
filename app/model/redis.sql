-- 公共说明
  聊天室ID = 唯一ID + 聊天室类型


-- 说明: 聊天室成员集合
-- 作用: 通过聊天室ID获取聊天室成员
set_im_chatroom_member:{$chatroom_id(聊天室ID)}
    value: {$member_id(会员ID)}

-- 说明: 聊天室消息记录
sorted_set_im_chatroom_message_record:{$chatroom_id(聊天室ID)}
    key: (系统消息ID)
    value: (消息string, 必须包含: 消息内容, 消息类型, 发送用户)

-- 说明: 储存聊天室消息ID
-- 作用:
hash_im_chatroom_message_id
    key: (聊天室ID)
    value: (消息ID , 每收到一条消息,消息ID加1)

-- 废弃
-- 说明: 储存聊天室 系统消息ID 跟 客户端消息ID关系 -
-- 作用:通过客户端消息ID, 可以查询到系统消息ID, 做消息同步|消息回执等等 需要使用
sorted_set_im_chatroom_client_message_id_join_server_message_id:{$chatroom_id(聊天室ID)}
    key: (系统消息ID)
    value: (客户端消息ID , 随机字符串,必须唯一)

-- 说明: 单聊业务
-- 用途: 通过自己的会员ID跟对方的会员ID获取聊天室ID
hash_im_chatroom_member_id_get_chatroom_id
    key: (用户ID从小到大排序, 并逗号分隔)
    value: (聊天室ID)

-- 说明: 群聊业务
-- 作用:
hash_
    key: (聊天室ID)
    value: (聊天室类型)


-- 说明: 消息送达失败
-- 作用: 消息发送失败, 下次设备上线自动发送
list_message_send_failure:{用户ID}
    value: (消息string)