package service

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"log"
	"ops-server/entity"
)

type Ldap struct {
	client *ldap.Conn
}

var DefaultLdap *Ldap

func CreateLdapConnection() (*Ldap, error) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "121.5.49.41", 389))
	if err != nil {
		log.Printf("ldap.Dial err is %s\n", err)
		return nil, err
	}

	DefaultLdap = &Ldap{client: l}
	if err := DefaultLdap.client.Bind("cn=admin,dc=langzhihe,dc=com", "xixianbin520"); err != nil {
		log.Printf("DefaultLdap.client.Bind err is %s\n", err)
		return nil, err
	}
	return DefaultLdap, nil

}

func (l *Ldap) CreateUser(p *entity.CreateUserParams) error {

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
		}}

	addRequest.Attributes = attr
	err := l.client.Add(addRequest)

	if err != nil {
		log.Printf("client.Add(addRequest) is err %s\n", err)
		return err
	}

	return nil

}

func (l *Ldap) DeleteUser(p *entity.DeleteUserParams) error {

	delRequest := ldap.NewDelRequest(p.Dn, nil)
	err := l.client.Del(delRequest)
	if err != nil {
		return err
	}

	return nil
}
