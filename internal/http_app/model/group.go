package model

// 群组表
type Group struct {
	GroupId         int64  `gorm:"column:group_id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"group_id"`
	Name            string `gorm:"column:name;type:char(50);NOT NULL" json:"name"`                                       // 群组名称
	Avatar          string `gorm:"column:avatar;type:char(50);NOT NULL" json:"avatar"`                                   // 群组头像
	Id              string `gorm:"column:id;type:varchar(20);NOT NULL" json:"id"`                                        // ID, 对用户展示并且唯一
	ChatroomId      string `gorm:"column:chatroom_id;type:char(32);NOT NULL" json:"chatroom_id"`                         // 房间ID
	OwnerMemberId   int64  `gorm:"column:owner_member_id;type:bigint(20);NOT NULL" json:"owner_member_id"`               // 所属者会员ID
	FounderMemberId int64  `gorm:"column:founder_member_id;type:bigint(20);default:0;NOT NULL" json:"founder_member_id"` // 创始人ID
	Permissions     string `gorm:"column:permissions;type:char(10);default:public;NOT NULL" json:"permissions"`          // 聊天室权限。 public:开放, protected:受保护(可见,并且管理员同意才能加入), private:私有(不可申请,并且管理员邀请才能加入)
	CreatedAt       int    `gorm:"column:created_at;type:int(11);default:0;NOT NULL" json:"created_at"`                  // 添加时间
}

func (m *Group) TableName() string {
	return "group"
}
