// Code generated by sql2gorm. DO NOT EDIT.
package model

// 动态表
type Dynamic struct {
	DynamicId        int64  `gorm:"column:dynamic_id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"dynamic_id"`
	MemberId         int64  `gorm:"column:member_id;type:bigint(20);NOT NULL" json:"member_id"`                          // 会员ID
	Content          string `gorm:"column:content;type:varchar(500);NOT NULL" json:"content"`                            // 内容
	Type             string `gorm:"column:type;type:char(10);default:common;NOT NULL" json:"type"`                       // 类型, 普通(文字或加图片):common, 视频:video
	ImageUrl         string `gorm:"column:image_url;type:varchar(1000);NOT NULL" json:"image_url"`                       // 图片地址
	VideoUrl         string `gorm:"column:video_url;type:varchar(255);NOT NULL" json:"video_url"`                        // 视频地址
	VideoCover       string `gorm:"column:video_cover;type:varchar(255);NOT NULL" json:"video_cover"`                    // 视频封面图
	VideoCoverWidth  int    `gorm:"column:video_cover_width;type:int(11);default:0;NOT NULL" json:"video_cover_width"`   // 视频封面图宽
	VideoCoverHeight int    `gorm:"column:video_cover_height;type:int(11);default:0;NOT NULL" json:"video_cover_height"` // 视频封面图高
	Zan              int    `gorm:"column:zan;type:int(11);default:0;NOT NULL" json:"zan"`                               // 点赞数
	Comment          int    `gorm:"column:comment;type:int(11);default:0;NOT NULL" json:"comment"`                       // 评论数
	AddressName      string `gorm:"column:address_name;type:varchar(50);NOT NULL" json:"address_name"`                   // 地址名称
	Latitude         string `gorm:"column:latitude;type:varchar(255);NOT NULL" json:"latitude"`                          // 经纬度: 经度
	Longitude        string `gorm:"column:longitude;type:varchar(255);NOT NULL" json:"longitude"`                        // 经纬度: 维度
	Purview          string `gorm:"column:purview;type:char(10);default:public;NOT NULL" json:"purview"`                 // 公开权限: public-公开, protected-好友可见, private-仅自己和指定用户可见
	PrivateToUid     string `gorm:"column:private_to_uid;type:text" json:"private_to_uid"`                               // 私有可见用户 逗号分隔
	Review           string `gorm:"column:review;type:char(10);default:wait;NOT NULL" json:"review"`                     // 审核状态: wait-审核中, normal-正常, refuse-拒绝
	DeletedAt        int    `gorm:"column:deleted_at;type:int(11);default:0;NOT NULL" json:"deleted_at"`                 // 删除时间
	CreatedAt        int    `gorm:"column:created_at;type:int(11);default:0;NOT NULL" json:"created_at"`                 // 添加时间
}

func (m *Dynamic) TableName() string {
	return "dynamic"
}
