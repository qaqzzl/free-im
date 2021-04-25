package model

// 用户会员
type UserMember struct {
	MemberId  uint   `gorm:"column:member_id;primary_key;AUTO_INCREMENT" json:"member_id"`
	Nickname  string `gorm:"column:nickname;NOT NULL" json:"nickname"`               // 用户昵称
	Id        string `gorm:"column:id;NOT NULL" json:"id"`                           // ID, 对用户展示并且唯一
	Gender    string `gorm:"column:gender;default:wz;NOT NULL" json:"gender"`        // wz-未知, w-女, m-男, z-中性
	Birthdate int    `gorm:"column:birthdate;default:0;NOT NULL" json:"birthdate"`   // 出生日期
	Avatar    string `gorm:"column:avatar;NOT NULL" json:"avatar"`                   // 头像
	Signature string `gorm:"column:signature;NOT NULL" json:"signature"`             // 个性签名
	City      string `gorm:"column:city;NOT NULL" json:"city"`                       // 城市
	Province  string `gorm:"column:province;NOT NULL" json:"province"`               // 省份
	CreatedAt int    `gorm:"column:created_at;default:0;NOT NULL" json:"created_at"` // 添加时间
	UpdatedAt int    `gorm:"column:updated_at;default:0;NOT NULL" json:"updated_at"` // 修改时间
	DeletedAt int    `gorm:"column:deleted_at;default:0;NOT NULL" json:"deleted_at"` // 删除时间
}

func (m *UserMember) TableName() string {
	return "user_member"
}
