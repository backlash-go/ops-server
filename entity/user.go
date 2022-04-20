package entity

import "time"

type RoleId struct {
	RoleId uint64
}

type UserIDRoleContact struct {
	UserId uint64 `json:"id"`
	Role   string `json:"role"`
}

type UserInfo struct {
	Token    string   `json:"token"`
	UserId   string   `json:"user_id"`
	UserName string   `json:"user_name"`
	Email    string   `json:"email"`
	Role     []string `json:"role"`
}

//分页查询用户
type UserList struct {
	Id          uint64    `json:"id"`
	UserName    string    `json:"user_name"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserInfoListRequest struct {
	PageSize   int    `json:"page_size" query:"page_size"` // 每页查询的数量
	Page       int    `json:"page" query:"page"`           // 查询的页码
	SearchName string `json:"search_name" query:"search_name"`
}

type UserInfoListResponse struct {
	Items      []UserList `json:"items"`
	TotalCount int64      `json:"total_count"`
}
