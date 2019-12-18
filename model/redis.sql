-- 说明: 聊天室成员集合
-- 作用: 通过聊天室ID获取聊天室成员
set_im_chatroom_member_{$chatroom_id(聊天室ID)}
    value: {$member_id(会员ID)}

-- 说明: 聊天室消息记录
hash_im_chatroom_message_record_{$chatroom_id(聊天室ID)}
    key: (消息ID)
    value: (消息string)

-- 说明: 单聊业务
-- 用途: 通过自己的会员ID跟对方的会员ID获取聊天室ID
hash_im_chatroom_member_id_get_chatroom_id
    key: (用户ID从小到大排序, 并逗号分隔)
    value: (聊天室ID)