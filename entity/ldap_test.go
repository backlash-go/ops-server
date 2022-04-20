package entity

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"log"
	"testing"
)

var L *ldap.Conn

func TestMain(m *testing.M) {
	log.Println("begin")

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "121.5.49.41", 389))
	if err != nil {
		log.Println("ldap.Dial err")

	}

	if err := l.Bind("cn=admin,dc=langzhihe,dc=com", "xixianbin520"); err != nil {
		panic(err)
		return
	}

	L = l
	log.Println("ldap succ")

	m.Run()
}

func TestLdap_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		args    *CreateUserParams
		wantErr bool
	}{
		{
			name: "createUser",
			args: &CreateUserParams{
				Cn:           "test11",
				Sn:           "test11",
				Mail:         "test11",
				GivenName:    "test11",
				EmployeeType: []string{"techs", "test11"},
				DisplayName:  "test11",
				UserPassword: "test11",
			},
			wantErr: false,
		},
	}
	l := &Ldap{
		Client: L,
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			if err := l.CreateUser(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLdap_SearchUser(t *testing.T) {

	tests := []struct {
		name    string
		args    *AuthUserParams
		wantErr bool
	}{
		{
			name: "searchUser",
			args: &AuthUserParams{UserPassword: "xixianbin205", Cn: "xixb"},

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Ldap{
				Client: L,
			}
			got, err := l.SearchUser(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got.Entries[0].Attributes)

 			for k, v := range got.Entries[0].Attributes {
				fmt.Println(k)

				fmt.Println(v.Values)
			}
		})
	}
}

func TestLdap_ModifyUserPassword(t *testing.T) {
	type args struct {
		newPass string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestLdap_ModifyUserPassword",
			args:args{newPass:"xixianbin222"},
			wantErr:false,

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Ldap{
				Client: L,
			}
			if err := l.ModifyUserPassword(tt.args.newPass); (err != nil) != tt.wantErr {
				t.Errorf("ModifyUserPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}