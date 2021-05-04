package model

// 会员授权账号表
type UserAuths struct {
	Id           int64  `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	MemberId     int64  `gorm:"column:member_id;type:bigint(20);NOT NULL" json:"member_id"`       // 会员ID
	IdentityType string `gorm:"column:identity_type;type:char(20);NOT NULL" json:"identity_type"` // 类型,wechat_applet,qq,wb,phone,number,email
	Identifier   string `gorm:"column:identifier;type:varchar(64);NOT NULL" json:"identifier"`    // 微信,QQ,微博openid | 手机号,邮箱,账号
	Credential   string `gorm:"column:credential;type:varchar(64);NOT NULL" json:"credential"`    // 密码凭证（站内的保存密码，站外的不保存或保存access_token）
}

func (m *UserAuths) TableName() string {
	return "user_auths"
}
