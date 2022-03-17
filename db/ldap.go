package db

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"log"
	"ops-server/entity"
)

var DefaultLdap *entity.Ldap

func InitLdap() {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "121.5.49.41", 389))
	if err != nil {
		log.Printf("ldap.Dial err is %s\n", err)
		panic(err)
		return
	}
	DefaultLdap = &entity.Ldap{Client: l}
	if err := DefaultLdap.Client.Bind("cn=admin,dc=langzhihe,dc=com", "xixianbin520"); err != nil {
		log.Printf("DefaultLdap.client.Bind err is %s\n", err)
		panic(err)

		return
	}
	log.Println("ldap connected")

	return

}

func GetLdap() *entity.Ldap {

	return DefaultLdap
}
