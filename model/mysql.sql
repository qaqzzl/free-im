

-- 聊天室表
create table if not exists `im_chatroom`(
    `chatroom_id` int(11) unsigned auto_increment primary key,
    `chatroom_class` char(10) not null comment '聊天室类型。 single_chat: 单聊, group_chat: 群聊',
    `recent_msg_time` int not null default 0 comment '最近消息时间',
    `created_at` int not null comment '创建时间'
)engine=innodb default charset=utf8 comment '聊天室表';

-- 聊天室成员表
create table if not exists `im_chatroom_member`(
    `chatroom_member_id` int(11) unsigned auto_increment primary key,
    `chatroom_id` int(11) not null comment '聊天室ID',
    `member_id` int(11) not null comment '会员ID',
    KEY `chatroom_id` (`chatroom_id`),
    KEY `member_id` (`member_id`)
)engine=innodb default charset=utf8 comment '聊天室成员表';

-- 聊天室群聊信息表
create table if not exists `im_chatroom_group_info`(
    `chatroom_id` int(11) unsigned primary key,
    `belong_member_id` int not null default 0 comment '所属会员ID',
    `founder_member_id` int not null default 0 comment '创始人ID',
    `chatroom_name` char(20) not null default '' comment '聊天室名称',
    `permissions` char(10) not null default '聊天室权限。 public:开放, protected:受保护(可见,并且管理员同意才能加入), private:私有(不可见,并且管理员邀请才能加入)',
    KEY `belong_member_id` (`belong_member_id`),
    KEY `founder_member_id` (`founder_member_id`)
)engine=innodb default charset=utf8 comment '聊天室表';

-- 聊天室消息记录
