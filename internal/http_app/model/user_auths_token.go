package model

// 用户授权 token 表
type UserAuthsToken struct {
	Id        int64  `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	MemberId  int64  `gorm:"column:member_id;type:bigint(20);NOT NULL" json:"member_id"`          // 会员ID
	Token     string `gorm:"column:token;type:varchar(255);NOT NULL" json:"token"`                // token
	Client    string `gorm:"column:client;type:char(20);NOT NULL" json:"client"`                  // app,web,wechat_applet
	LastTime  int    `gorm:"column:last_time;type:int(11);NOT NULL" json:"last_time"`             // 上次刷新时间
	Status    int    `gorm:"column:status;type:tinyint(1);default:0;NOT NULL" json:"status"`      // 1-其他设备强制下线
	CreatedAt int    `gorm:"column:created_at;type:int(11);default:0;NOT NULL" json:"created_at"` // 添加时间
}

func (m *UserAuthsToken) TableName() string {
	return "user_auths_token"
}
