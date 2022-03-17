package db

import (
	"flag"
	"github.com/spf13/viper"
	"log"
	"ops-server/utils"
)

func Init()  {
	configFile := flag.String("conf", "config/dev-config.yaml", "path of config file")
	flag.Parse()
	viper.SetConfigFile(*configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper read config is failed, err is %v configFile is %s \n", err, configFile)
	}
	log.Println("logger init ")
	//init mysql
	dbConf := viper.GetStringMapString("database")
	log.Println(dbConf)


	utils.InitLogger(viper.GetString("env"))

	InitDB(dbConf["user"], dbConf["password"], dbConf["host"], dbConf["port"], dbConf["name"])

	//dbRedisConf := viper.GetStringMapString("authRedis")
	//InitRedis(fmt.Sprintf("%s:%s", dbRedisConf["host"], dbRedisConf["port"]),dbRedisConf["password"],dbRedisConf["db"])
}
