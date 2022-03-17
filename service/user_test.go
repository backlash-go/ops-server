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
			args:    args{cn: "test1"},
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
