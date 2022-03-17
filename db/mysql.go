package db

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var db *gorm.DB

// InitDB 初始化 MySQL 链接
func InitDB(user, password, host, port, dbName string) {
	log.Println("connecting MySQL ... ", host)
	mdb, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, password, host, port, dbName))
	if err != nil {
		panic(err)
		return
	}
	if mdb == nil {
		panic("failed to connect database")
	}
	mdb.LogMode(true)

	if viper.GetString("env") == "online" {
		mdb.LogMode(false)
	}

	log.Println("connected")
	db = mdb
	return
}

// GetDB 获取数据库链接实例
func GetDB() *gorm.DB {
	return db
}
