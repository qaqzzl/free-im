package model

// 群组表
type Group struct {
	GroupId         int64  `gorm:"column:group_id;primary_key;AUTO_INCREMENT" json:"group_id"`
	Name            string `gorm:"column:name;NOT NULL" json:"name"`                                                   // 群组名称
	Avatar          string `gorm:"column:avatar;NOT NULL" json:"avatar"`                                               // 群组头像
	Desc            string `gorm:"column:desc;NOT NULL" json:"desc"`                                                   // 描述
	Id              string `gorm:"column:id;NOT NULL" json:"id"`                                                       // ID, 对用户展示并且唯一
	ChatroomId      int64  `gorm:"column:chatroom_id;NOT NULL" json:"chatroom_id"`                                     // 房间ID
	OwnerMemberId   int64  `gorm:"column:owner_member_id;NOT NULL" json:"owner_member_id"`                             // 所属者会员ID
	FounderMemberId int64  `gorm:"column:founder_member_id;default:0;NOT NULL" json:"founder_member_id"`               // 创始人ID
	Permissions     string `gorm:"column:permissions;default:public;NOT NULL" json:"permissions"`                      // 聊天室权限。 public:开放, protected:受保护(可见,并且管理员同意才能加入), private:私有(不可申请,并且管理员邀请才能加入)
	CreatedAt       int    `gorm:"autoCreateTime;column:created_at;type:int(11);default:0;NOT NULL" json:"created_at"` // 添加时间
	UpdatedAt       int    `gorm:"autoUpdateTime;column:updated_at;type:int(11);default:0;NOT NULL" json:"updated_at"` // 修改时间
}

func (m *Group) TableName() string {
	return "group"
}
