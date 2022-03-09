package entity

type CreateUserParams struct {
	Cn           string   `json:"cn"`
	Sn           string   `json:"sn"`
	Mail         string   `json:"mail"`
	GivenName    string   `json:"given_name"`
	EmployeeType []string `json:"employee_type"`
	DisplayName  string   `json:"display_name"`
	UserPassword string   `json:"user_password"`
}


type DeleteUserParams struct {
	Dn string `json:"dn"`
}
