package models

import "time"

type Permission struct {
	Id        uint64     `gorm:"column:id" form:"id" json:"id"`
	Api       string     `gorm:"column:user_name" form:"user_name" json:"user_name"`
	Name      string     `gorm:"column:name", form:"name",json:"name"`
	CreatedAt *time.Time `gorm:"column:created_at" form:"created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" form:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" form:"deleted_at" json:"deleted_at"`
}

//帐号密码登陆验证
func (m *Permission) TableName() string {
	return "Permission"

}
