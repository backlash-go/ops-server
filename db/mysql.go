package db

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var db *gorm.DB

// InitDB 初始化 MySQL 链接
func InitDB(user, password, host, port, dbName string) {


	loggerLevel := logger.Info
	if viper.GetString("env") == "online" {
		loggerLevel = logger.Silent
	}


	mdb, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)), &gorm.Config{
		//SkipDefaultTransaction: true,
		//PrepareStmt:            true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				//SlowThreshold: time.Second,   // Slow SQL threshold
				LogLevel: loggerLevel, // Log level
			},)})

	log.Println("connecting MySQL ... ", host)

	if err != nil {
		panic(err)
		return
	}
	if mdb == nil {
		panic("failed to connect database")
	}

	log.Println("connected")
	db = mdb
	return
}

// GetDB 获取数据库链接实例
func GetDB() *gorm.DB {
	return db
}
