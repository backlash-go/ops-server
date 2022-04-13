package service

import (
	"fmt"
	"ops-server/models"
	"testing"
)

func TestQueryUser(t *testing.T) {
	type args struct {
		cn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestQueryUser",
			args:    args{cn: "test11"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := QueryUser(tt.args.cn)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("%+v", gotUser)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type args struct {
		userId  uint64
		updates interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{
			name: "TestQueryUser",
			args: args{userId: 1, updates: map[string]interface{}{
				"email": "xixianbin11@qq.com",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateUser(tt.args.userId, tt.args.updates); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestAddUser(t *testing.T) {
	type args struct {
		req models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestAddUser",
			args:    args{req: models.User{UserName: "test6", Email: "test6@qq.com"}},
			wantErr: false,
		},}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddUser(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Printf("id is : %+v", got)

		})
	}
}

func TestQueryRoleId(t *testing.T) {
	type args struct {
		userId uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestQueryRoleId",
			args:    args{userId: 12},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRoleIDs, err := QueryUserRoleId(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryRoleId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("id is : %+v", gotRoleIDs)

		})
	}
}

func TestQueryRole(t *testing.T) {
	type args struct {
		roleId []uint64
	}
	tests := []struct {
		name          string
		args          args
		wantRoleNames []string
		wantErr       bool
	}{
		{
			name:    "TestQueryRoleId",
			args:    args{roleId: []uint64{2, 3}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRoleNames, err := QueryUserRoles(tt.args.roleId)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("role_name is : %+v", gotRoleNames)

		})
	}
}

func TestQueryAllUser(t *testing.T) {
	type args struct {
		userName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestQueryAllUser",
			args: args{userName: "sad"},
		},}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := QueryAllUser(tt.args.userName)
			if err != nil {
				t.Errorf("QueryAllUser() error = %v", err)
				return
			}
			fmt.Printf("role_name is : %+v", gotUser)

		})
	}
}

func TestQueryUserRoleIdByRoleName(t *testing.T) {
	type args struct {
		roleName []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestQueryUserRoleIdByRoleName",
			args:    args{roleName: []string{"admin", "devops"},},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRoleIDs, err := QueryUserRoleIdByRoleName(tt.args.roleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryUserRoleIdByRoleName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("role_id is : %+v", gotRoleIDs)

		})
	}
}

func TestQueryPermissionIdByRoleId(t *testing.T) {
	type args struct {
		roleId []uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{{name: "TestQueryPermissionIdByRoleId", args: args{roleId: []uint64{11, 11, 12}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPermissionId, err := QueryPermissionIdByRoleId(tt.args.roleId)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryPermissionIdByRoleId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("role_id is : %+v", gotPermissionId)

		})
	}
}

func TestQueryApi(t *testing.T) {
	type args struct {
		permissionId []uint64
	}
	tests := []struct {
		name     string
		args     args
		wantApis []string
		wantErr  bool
	}{
		{name: "TestQueryApi", args: args{permissionId: []uint64{1, 2}},wantErr:false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotApis, err := QueryApi(tt.args.permissionId)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryApi() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("gotApis is : %+v", gotApis)

		})
	}
}

func TestCreateUserRecord(t *testing.T) {
	type args struct {
		user models.User
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{name: "TestCreateUserRecord", args: args{user:models.User{UserName:"test98888",Email:"test98888@qq.com"}},wantErr:false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateUserRecord(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUserRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("got is : %+v", got)

		})
	}
}

func TestCreateUserRoleRecord(t *testing.T) {
	type args struct {
		userRole models.UserRole
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestCreateUserRoleRecord", args: args{userRole:models.UserRole{UserId:uint64(8),RoleId:uint64(3)}},wantErr:false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateUserRoleRecord(tt.args.userRole); (err != nil) != tt.wantErr {
				t.Errorf("CreateUserRoleRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}



func TestAddUserRoles1(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "TestAddUserRoles1",wantErr:false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddUserRoles(); (err != nil) != tt.wantErr {
				t.Errorf("AddUserRoles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}