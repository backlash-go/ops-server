package db

import (
	"log"
	"testing"
)

func TestMain(m *testing.M) {

	log.Println("11")

	InitRedis("127.0.0.1","6379","sahj2314sahdj","11")


	m.Run()
}
