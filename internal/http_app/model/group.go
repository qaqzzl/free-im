package model

type Group struct {
	GroupId         uint `gorm:"primaryKey"`
	Id              string
	Name            string
	Avatar          string
	ChatroomId      string
	OwnerMemberId   string
	FounderMemberId string
	Permissions     string
	CreatedAt       int
}
