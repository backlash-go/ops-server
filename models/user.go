package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id          uint64         `gorm:"column:id" form:"id" json:"id"`
	UserName    string         `gorm:"column:user_name" form:"user_name" json:"user_name"`
	Email       string         `gorm:"column:email", form:"email",json:"email"`
	DisplayName string         `gorm:"column:display_name", form:"display_name",json:"display_name"`
	CreatedAt   time.Time      `gorm:"column:created_at" form:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" form:"updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" form:"deleted_at" json:"deleted_at"`
}

//帐号密码登陆验证
func (m *User) TableName() string {
	return "user"

}
