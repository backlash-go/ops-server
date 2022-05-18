package entity

import (
	"time"
)

type PermissionInfoRequest struct {
	ID uint64 `json:"id" query:"id"`
}

type UpdateApiParamsRequest struct {
	ID   uint64   `json:"id"`
	Api  string   `json:"api`
	Name string   `json:"name"`
	Role []string `json:"role"`
}

type CreateApiParamsRequest struct {
	Api  string   `json:"api`
	Name string   `json:"name"`
	Role []uint64 `json:"role"`
}

type PermissionInfoList struct {
	Id        uint64    `json:"id"`
	Api       string    `json:"api"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type PermissionInfoListRequest struct {
	PageSize   int    `json:"page_size" query:"page_size"`     // 每页查询的数量
	Page       int    `json:"page" query:"page"`               // 查询的页码
	SearchName string `json:"search_name" query:"search_name"` //
}

type PermissionInfoListResponse struct {
	Items      []PermissionInfoList `json:"items"`
	TotalCount int64                `json:"total_count"`
}

type PermissionIdRoleContact struct {
	PermissionId uint64 `json:"permission_id"`
	Role         string `json:"role"`
}
