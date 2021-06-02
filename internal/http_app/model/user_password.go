package model

// 会员账号密码表
type UserPassword struct {
	Id        uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	MemberId  int64  `gorm:"column:member_id;NOT NULL"`            // 会员ID
	Pwd       string `gorm:"column:pwd;NOT NULL"`                  // 密码，加密后
	Status    int    `gorm:"column:status;default:1;NOT NULL"`     // 状态，0：正常，1：失效，2：禁用
	CreatedAt int    `gorm:"column:created_at;default:0;NOT NULL"` // 添加时间
	UpdatedAt int    `gorm:"column:updated_at;default:0;NOT NULL"` // 修改时间
}

func (m *UserPassword) TableName() string {
	return "user_password"
}
