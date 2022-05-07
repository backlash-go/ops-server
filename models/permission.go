package models

import (
	"gorm.io/gorm"
	"time"
)

type Permission struct {
	Id        uint64         `gorm:"column:id" form:"id" json:"id"`
	Api       string         `gorm:"column:api" form:"api" json:"api"`
	Name      string         `gorm:"column:name", form:"name",json:"name"`
	CreatedAt time.Time      `gorm:"column:created_at" form:"created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" form:"updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" form:"deleted_at" json:"deleted_at"`
}

//帐号密码登陆验证
func (m *Permission) TableName() string {
	return "Permission"

}
