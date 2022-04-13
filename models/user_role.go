package models

import "time"

type UserRole struct {
	Id        uint64     `gorm:"column:id" form:"id" json:"id"`
	UserId    uint64     `gorm:"column:user_id" form:"user_id" json:"user_id"`
	RoleId    uint64   `gorm:"column:role_id" form:"role_id" json:"role_id"`
	CreatedAt *time.Time `gorm:"column:created_at" form:"created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" form:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" form:"deleted_at" json:"deleted_at"`
}

//帐号密码登陆验证
func (m *UserRole) TableName() string {
	return "user_role"

}
