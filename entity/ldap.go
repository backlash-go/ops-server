package entity

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"log"
)

type CreateUserParams struct {
	Cn           string   `json:"cn"`
	Sn           string   `json:"sn"`
	Mail         string   `json:"mail"`
	GivenName    string   `json:"given_name"`
	EmployeeType []string `json:"employee_type"`
	DisplayName  string   `json:"display_name"`
	UserPassword string   `json:"user_password"`
	Role         []uint64 `json:"role"`
}

type ModifyUserPassword struct {
	Password string `json:"password"`
	Cn       string `json:"cn"`
}

type DeleteUserParams struct {
	Cn string `json:"cn"`
}

type AuthUserParams struct {
	Cn           string `json:"cn"`
	UserPassword string `json:"user_password"`
}

type LoginAuthResp struct {
	Token        string   `json:"token"`
	UserId       uint64   `json:"user_id"`
	UserName     string   `json:"user_name"`
	Email        string   `json:"email"`
	Role         []string `json:"role"`
	EmployeeType string   `json:"employee_type"`
	DisplayName  string   `json:"display_name"`
}

type LdapUserInfo struct {
	Cn           string   `json:"cn" query:"cn"`
	Mail         string   `json:"mail" query:"mail"`
	DisPlayName  string   `json:"display_name" query:"display_name"`
	EmployeeType []string `json:"employee_type" query:"employee_type"`
	Role         []string `json:"role" query:"role"`
}

type Ldap struct {
	Client *ldap.Conn
}

func (l *Ldap) SearchUser(u *AuthUserParams) (*ldap.SearchResult, error) {

	searchRequest := ldap.NewSearchRequest(
		"dc=langzhihe,dc=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(userPassword=%s)(cn=%s))", u.UserPassword, u.Cn),
		[]string{"dn", "mail", "sn", "givenName", "cn", "displayName", "employeeType", "objectClass", "userPassword"},
		nil,
	)

	result, err := l.Client.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (l *Ldap) CreateUser(p *CreateUserParams) error {

	addRequest := ldap.NewAddRequest(fmt.Sprintf("cn=%s,ou=person,dc=langzhihe,dc=com", p.Cn))
	var attr = []ldap.Attribute{
		{
			Type: "objectClass",
			Vals: []string{"inetOrgPerson"},
		},
		{
			Type: "employeeType",
			Vals: p.EmployeeType,
		},
		{
			Type: "cn",
			Vals: []string{p.Cn},
		}, {
			Type: "sn",
			Vals: []string{p.Sn},
		}, {
			Type: "uid",
			Vals: []string{p.Cn},
		}, {
			Type: "givenName",
			Vals: []string{p.GivenName},
		}, {
			Type: "userPassword",
			Vals: []string{p.UserPassword},
		},
		{
			Type: "title",
			Vals: []string{fmt.Sprintf("%s-title", p.Cn)},
		},
		{
			Type: "mail",
			Vals: []string{p.Mail},
		},
		{
			Type: "displayName",
			Vals: []string{p.DisplayName},
		}}

	addRequest.Attributes = attr
	err := l.Client.Add(addRequest)

	if err != nil {
		log.Printf("Client.Add(addRequest) is err %s\n", err)
		return err
	}

	return nil

}

func (l *Ldap) DeleteUser(p string) error {

	delRequest := ldap.NewDelRequest(p, nil)
	err := l.Client.Del(delRequest)
	if err != nil {
		return err
	}

	return nil
}

func (l *Ldap) ModifyUserPassword(req *ModifyUserPassword) error {

	dn := fmt.Sprintf("cn=%s,ou=person,dc=langzhihe,dc=com", req.Cn)

	modifyRequest := ldap.NewModifyRequest(dn)
	modifyRequest.Replace("userPassword", []string{req.Password})

	err := l.Client.Modify(modifyRequest)

	if err != nil {
		return err
	}
	return nil

}

func (l *Ldap) ModifyUserInfo(req *LdapUserInfo) error {

	dn := fmt.Sprintf("cn=%s,ou=person,dc=langzhihe,dc=com", req.Cn)

	modifyRequest := ldap.NewModifyRequest(dn)

	modifyRequest.ReplaceAttributes = []ldap.PartialAttribute{{Type: "mail", Vals: append(ldap.PartialAttribute{}.Vals, req.Mail)},
		{Type: "displayName", Vals: append(ldap.PartialAttribute{}.Vals, req.DisPlayName)}, {Type: "employeeType", Vals: req.EmployeeType},}
	err := l.Client.Modify(modifyRequest)

	if err != nil {
		return err
	}
	return nil

}
