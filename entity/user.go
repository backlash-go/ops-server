package entity

type RoleId struct {
	RoleId uint64
}


type UserInfo struct {
	Token    string `json:"token"`
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Role     []string `json:"role"`
}
