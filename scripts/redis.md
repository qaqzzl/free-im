-- 公共说明



-- 说明: 聊天室成员集合
-- 作用: 通过聊天室ID获取聊天室成员
set_im_chatroom_member:{$chatroom_id(聊天室ID)}
    value: {$member_id(会员ID)}

-- 说明: 消息记录(写扩散)
sorted_set_im_user_message_record:{$user_id(用户ID)}
    key: (decode(消息ID))
    value: (消息)

-- 说明: 消息记录(读扩散)
sorted_set_im_chatroom_message_record:{$chatroom_id(聊天室ID)}
    key: (decode(消息ID))
    value: (消息)

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


-- 说明: 消息离线
-- 作用: 消息离线, 下次设备上线自动发送
list_message_offline:{用户ID}
    value: (消息)

-- ack超时 消息重传
list_message_ack_timeout_retransmit:{用户ID}:(客户端类型)
    value: (消息ID)
hash_message_ack_timeout_retransmit:{用户ID}:(客户端类型)
    key: (消息ID)
    value: (重传时间, 毫秒) | (消息)