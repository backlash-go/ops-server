package service

import (
	"ops-server/db"
	"testing"
)

func TestMain(m *testing.M) {

	db.InitDB("devops","fvf7sfRj23ns1adf","127.0.0.1","3306","ops")


	m.Run()
}
