/*
 Navicat Premium Data Transfer

 Source Server         : free-im
 Source Server Type    : MySQL
 Source Server Version : 100412
 Source Host           : 101.132.107.212:3306
 Source Schema         : free_im

 Target Server Type    : MySQL
 Target Server Version : 100412
 File Encoding         : 65001

 Date: 30/04/2021 11:37:54
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for dynamic
-- ----------------------------
DROP TABLE IF EXISTS `dynamic`;
CREATE TABLE `dynamic` (
  `dynamic_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `member_id` bigint(20) NOT NULL COMMENT '会员ID',
  `content` varchar(500) NOT NULL DEFAULT '' COMMENT '内容',
  `type` char(10) NOT NULL DEFAULT 'common' COMMENT '类型, 普通(文字或加图片):common, 视频:video',
  `image_url` varchar(1000) NOT NULL DEFAULT '' COMMENT '图片地址',
  `video_url` varchar(255) NOT NULL DEFAULT '' COMMENT '视频地址',
  `video_cover` varchar(255) NOT NULL DEFAULT '' COMMENT '视频封面图',
  `video_cover_width` int(11) NOT NULL DEFAULT 0 COMMENT '视频封面图宽',
  `video_cover_height` int(11) NOT NULL DEFAULT 0 COMMENT '视频封面图高',
  `zan` int(11) NOT NULL DEFAULT 0 COMMENT '点赞数',
  `comment` int(11) NOT NULL DEFAULT 0 COMMENT '评论数',
  `address_name` varchar(50) NOT NULL DEFAULT '' COMMENT '地址名称',
  `latitude` varchar(255) NOT NULL DEFAULT '' COMMENT '经纬度: 经度',
  `longitude` varchar(255) NOT NULL DEFAULT '' COMMENT '经纬度: 维度',
  `purview` char(10) NOT NULL DEFAULT 'public' COMMENT '公开权限: public-公开, protected-好友可见, private-仅自己和指定用户可见',
  `private_to_uid` text DEFAULT NULL COMMENT '私有可见用户 逗号分隔',
  `review` char(10) NOT NULL DEFAULT 'wait' COMMENT '审核状态: wait-审核中, normal-正常, refuse-拒绝',
  `deleted_at` int(11) NOT NULL DEFAULT 0 COMMENT '删除时间',
  `created_at` int(11) NOT NULL DEFAULT 0 COMMENT '添加时间',
  PRIMARY KEY (`dynamic_id`),
  KEY `member_id` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='动态表';

-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group` (
  `group_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` char(50) NOT NULL COMMENT '群组名称',
  `avatar` char(255) NOT NULL COMMENT '群组头像',
  `desc` varchar(255) not null default '' comment '描述',
  `id` varchar(20) NOT NULL COMMENT 'ID, 对用户展示并且唯一',
  `chatroom_id` bigint(20) NOT NULL COMMENT '房间ID',
  `owner_member_id` bigint(20) NOT NULL COMMENT '所属者会员ID',
  `founder_member_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '创始人ID',
  `permissions` char(10) NOT NULL DEFAULT 'public' COMMENT '聊天室权限。 public:开放, protected:受保护(可见,并且管理员同意才能加入), private:私有(不可申请,并且管理员邀请才能加入)',
  `created_at` int(11) NOT NULL DEFAULT 0 COMMENT '添加时间',
  `updated_at` int(11) NOT NULL DEFAULT 0 COMMENT '修改时间',
  PRIMARY KEY (`group_id`),
  UNIQUE KEY `id` (`id`),
  KEY `owner_member_id` (`owner_member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群组表';

-- ----------------------------
-- Table structure for group_member
-- ----------------------------
DROP TABLE IF EXISTS `group_member`;
CREATE TABLE `group_member` (
  `group_member_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `group_id` bigint(20) NOT NULL COMMENT '群组ID',
  `member_id` bigint(20) NOT NULL COMMENT '会员ID',
  `alias` varchar(50) NOT NULL COMMENT '会员群别名',
  `notify_level` tinyint(1) not null default 0 comment '通知级别，0：正常，1：接收消息但不提醒，2：屏蔽群消息',
  `member_identity` char(10) NOT NULL COMMENT '成员身份: admin-管理员, root-群主, common-普通成员',
  `status` char(10) NOT NULL COMMENT '状态: normal-正常, blacklist-黑名单',
  `created_at` int(11) NOT NULL DEFAULT 0 COMMENT '添加时间',
  PRIMARY KEY (`group_member_id`),
  KEY `group_id` (`group_id`),
  KEY `member_id` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群组成员表';

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `message_id` char(32) NOT NULL COMMENT '消息ID',
  `chatroom_id` bigint(20) NOT NULL COMMENT '聊天室ID',
  `member_id` bigint(20) NOT NULL COMMENT '发送消息会员ID',
  `content` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `chatroom_id` (`chatroom_id`),
  KEY `member_id` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息记录（储存）表';

-- 用户消息记录表，对于（单聊，普通群聊）储存用户消息记录
DROP TABLE IF EXISTS `user_message`;
CREATE TABLE `user_message` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `message_id` char(32) NOT NULL COMMENT '消息ID',
    `member_id` bigint(20) NOT NULL COMMENT '会员ID',
    PRIMARY KEY (`id`),
    KEY `message_id` (`message_id`),
    KEY `member_id` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户消息记录表';

DROP TABLE IF EXISTS `chatroom_record`;
CREATE TABLE `user_chatroom_record` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `member_id` bigint(20) NOT NULL COMMENT '会员ID',
    `chatroom_id` bigint(20) NOT NULL COMMENT '聊天室ID',
    `sort` char(32) NOT NULL COMMENT '排序',
    `expand` varchar(2000) DEFAULT NULL COMMENT '扩展',
    PRIMARY KEY (`id`),
    KEY `chatroom_id` (`chatroom_id`),
    KEY `member_id` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天对话记录表';

-- ----------------------------
-- Table structure for user_auths
-- ----------------------------
DROP TABLE IF EXISTS `user_auths`;
CREATE TABLE `user_auths` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `member_id` bigint(20) NOT NULL COMMENT '会员ID',
  `identity_type` char(20) NOT NULL COMMENT '类型,wechat_applet,qq,wb,phone,number,email',
  `identifier` varchar(64) NOT NULL DEFAULT '' COMMENT '微信,QQ,微博openid | 手机号,邮箱,账号',
  `credential` varchar(64) NOT NULL DEFAULT '' COMMENT '密码凭证（站外的不保存或保存access_token）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `identity_type_identifier` (`identity_type`,`identifier`) USING BTREE,
  KEY `member_id` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员授权账号表';

DROP TABLE IF EXISTS `user_password`;
CREATE TABLE `user_password` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `member_id` bigint(20) NOT NULL COMMENT '会员ID',
    `pwd` varchar(64) NOT NULL DEFAULT '' COMMENT '密码，加密后',
    `status` tinyint(1) not null default 1 COMMENT '状态，0：正常，1：失效，2：禁用',
    `created_at` int(11) NOT NULL DEFAULT 0 COMMENT '添加时间',
    `updated_at` int(11) NOT NULL DEFAULT 0 COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `member_id` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员账号密码表';

-- ----------------------------
-- Table structure for user_auths_token
-- ----------------------------
DROP TABLE IF EXISTS `user_auths_token`;
CREATE TABLE `user_auths_token` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `member_id` bigint(20) NOT NULL COMMENT '会员ID',
  `token` varchar(255) NOT NULL DEFAULT '' COMMENT 'token',
  `client` char(20) NOT NULL COMMENT 'app,web,wechat_applet',
  `last_time` int(11) NOT NULL COMMENT '上次刷新时间',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '1-其他设备强制下线',
  `created_at` int(11) NOT NULL DEFAULT 0 COMMENT '添加时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户授权 token 表';

-- ----------------------------
-- Table structure for user_friend
-- ----------------------------
DROP TABLE IF EXISTS `user_friend`;
CREATE TABLE `user_friend` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `member_id` bigint(20) NOT NULL COMMENT '会员ID',
  `friend_id` bigint(20) NOT NULL COMMENT '好友ID',
  `friend_remark` varchar(50) NOT NULL DEFAULT '' COMMENT '昵称备注',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0-正常, 1-删除',
  `created_at` int(11) NOT NULL DEFAULT 0 COMMENT '添加时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `member_id_friend_id` (`member_id`,`friend_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户好友表';

-- ----------------------------
-- Table structure for user_friend_apply
-- ----------------------------
DROP TABLE IF EXISTS `user_friend_apply`;
CREATE TABLE `user_friend_apply` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `member_id` bigint(20) NOT NULL COMMENT '会员ID',
  `friend_id` bigint(20) NOT NULL COMMENT '好友ID',
  `remark` varchar(50) NOT NULL DEFAULT '' COMMENT '添加好友备注',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0-等待, 1-同意, 2-拒绝',
  `created_at` int(11) NOT NULL DEFAULT 0 COMMENT '添加时间',
  PRIMARY KEY (`id`),
  KEY `member_id` (`member_id`),
  KEY `friend_id` (`friend_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='好友申请表';

-- ----------------------------
-- Table structure for user_member
-- ----------------------------
DROP TABLE IF EXISTS `user_member`;
CREATE TABLE `user_member` (
  `member_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `nickname` varchar(50) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `id` varchar(20) NOT NULL COMMENT 'ID, 对用户展示并且唯一',
  `gender` char(5) NOT NULL DEFAULT 'wz' COMMENT 'wz-未知, w-女, m-男, z-中性',
  `birthdate` int(11) NOT NULL DEFAULT 0 COMMENT '出生日期',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `signature` varchar(64) NOT NULL DEFAULT '' COMMENT '个性签名',
  `city` char(50) NOT NULL DEFAULT '' COMMENT '城市',
  `province` char(50) NOT NULL DEFAULT '' COMMENT '省份',
  `created_at` int(11) NOT NULL DEFAULT 0 COMMENT '添加时间',
  `updated_at` int(11) NOT NULL DEFAULT 0 COMMENT '修改时间',
  `deleted_at` int(11) NOT NULL DEFAULT 0 COMMENT '删除时间',
  PRIMARY KEY (`member_id`),
  UNIQUE KEY `nickname` (`nickname`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户会员';

SET FOREIGN_KEY_CHECKS = 1;
