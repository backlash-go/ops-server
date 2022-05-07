package models

import (
	"gorm.io/gorm"
	"time"
)

type RolePermission struct {
	Id           uint64         `gorm:"column:id" form:"id" json:"id"`
	RoleId       uint64          `gorm:"column:role_id" form:"role_id" json:"role_id"`
	PermissionId uint64          `gorm:"column:permission_id", form:"permission_id",json:"permission_id"`
	CreatedAt    time.Time      `gorm:"column:created_at" form:"created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" form:"updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" form:"deleted_at" json:"deleted_at"`
}

//帐号密码登陆验证
func (m *RolePermission) TableName() string {
	return "role_permission"

}
