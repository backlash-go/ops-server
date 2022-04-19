package service

import (
	"ops-server/db"
	"testing"
)

func TestMain(m *testing.M) {
	//db.Init()

	//configFile := flag.String("conf", "config/dev-config.yaml", "path of config file")
	//flag.Parse()
	//viper.SetConfigFile(*configFile)
	//if err := viper.ReadInConfig(); err != nil {
	//	log.Fatalf("viper read config is failed, err is %v configFile is %v \n", err, configFile)
	//}

	db.InitDB("devops","fvf7sfRj23ns1adf","127.0.0.1","3306","ops")


	m.Run()
}
