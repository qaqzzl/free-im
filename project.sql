-- 创建库
create database if not exists `free_im`;

use `free_im`;

-- 用户会员
create table if not exists `user_member`(
    `member_id` int unsigned auto_increment primary key,
    `nickname` varchar(50) not null default '' comment '用户昵称',
    `gender` char(5) not null default 'wz' comment 'wz-未知, w-女, m-男, z-中性',
    `birthdate` int not null default 0 comment '出生日期',
    `avatar` varchar(255) not null default '' comment '头像',
    `signature` varchar(64) not null default '' comment '个性签名',
    `city` char(50) not null default '' comment '城市',
    `province` char(50) not null default '' comment '省份',
    `created_at` int not null default 0 comment '添加时间',
    `updated_at` int not null default 0 comment '修改时间',
    `deleted_at` int not null default 0 comment '删除时间',
    UNIQUE KEY `nickname` (`nickname`)
)engine=innodb default charset=utf8mb4 comment '用户会员';

-- 会员授权账号表
create table if not exists `user_auths`(
    `id` int unsigned auto_increment primary key,
    `member_id` int not null comment '会员ID',
    `identity_type` char(20) not null comment '类型,wechat_applet,qq,wb,phone,number,email',
    `identifier` varchar(64) not null default '' comment '微信,QQ,微博opendid | 手机号,邮箱,账号',
    `credential` varchar(64) not null default '' comment '密码凭证（站内的保存密码，站外的不保存或保存access_token）',
    KEY `member_id` (`member_id`),
    UNIQUE KEY `identity_type_identifier` (`identity_type`,`identifier`) USING BTREE
)engine=innodb default charset=utf8 comment '会员授权账号表';

-- 用户授权 token 表 ,这个表用redis比较好 , 也可以使用JWS
create table if not exists `user_auths_token`(
    `id` int unsigned auto_increment primary key,
    `member_id` int not null comment '会员ID',
    `token` varchar(255) not null default '' comment 'token',
    `client` char(20) not null comment 'app,web,wechat_applet',
    `last_time` int not null comment '上次刷新时间',
    `status` tinyint(1) not null default 0 comment '1-其他设备强制下线',
    `created_at` int not null default 0 comment '添加时间',
    UNIQUE KEY `token` (`token`)
)engine=innodb default charset=utf8 comment '用户授权 token 表';


-- 用户好友表
create table if not exists `user_friend`(
    `id` int unsigned auto_increment primary key,
    `member_id` int not null comment '会员ID',
    `friend_id` int not null comment '好友ID',
    `member_remark` varchar(50) not null default "" comment '会员备注',
    `friend_remark` varchar(50) not null default "" comment '好友备注',
    `status` tinyint(1) not null default 0 comment '0-等待同意, 1-正常',
    `created_at` int not null default 0 comment '添加时间',
    KEY `member_id` (`member_id`),
    KEY `friend_id` (`friend_id`)
)engine=innodb default charset=utf8 comment '用户好友表';
-- 查询我的好友
select * from `user_friend` where `member_id` = {$uid} or `friend_id` = {$uid}