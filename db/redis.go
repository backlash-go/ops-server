package db

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

// RDClient redis 的客户端
var RDClient *redis.Client

// InitRedis 初始化redis
func InitRedis(host, port, password, db string) {

	dsn := strings.Join([]string{host, port}, ":")
	Db, _ := strconv.Atoi(db)
	RDClient = redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: password,
		DB:       Db,
	})
	pong, err := RDClient.Ping().Result()
	if err != nil {
		log.Println("failed to connect redis")
		panic(err)
	}

	log.Printf("redis connected ping is %s\n", pong)

	return
}

// GetRedis 获取redis链接实例
func GetRedis() *redis.Client {
	return RDClient
}

func RedisHMSet(token string, keyFields map[string]interface{}) error {

	return RDClient.HMSet(token, keyFields).Err()

}

func RedisHMGet(token string, keyFields ...string) ([]string, error) {
	tmpSlice := make([]string, 0, 0)
	res, err := RDClient.HMGet(token, keyFields...).Result()
	for k, _ := range res {
		if res[k] == nil {
			return nil, errors.New(fmt.Sprintf("redis HMGet key不存在 token is %s and keyFields is %s \n", token, keyFields))
		}
		v, ok := res[k].(string)
		if !ok {
			return nil, errors.New("RedisHMGet 断言失败")
		}
		tmpSlice = append(tmpSlice, v)
	}
	return tmpSlice, err

}

func RedisHGetAll(token string) (map[string]string, error) {

	redisMap, err := RDClient.HGetAll(token).Result()
	if err != nil {
		return nil, err
	}
	if len(redisMap) == 0 {
		return nil, nil
	}
	return redisMap, nil

}

func RedisSetKeyTtl(token string, expire time.Duration) error {
	return RDClient.Expire(token, expire).Err()
}

func RedisKeyIsExist(token string) (int64, error) {
	res, err := RDClient.Exists(token).Result()

	if err != nil {
		return 0, err
	}
	return res, err
}

func RedisDelKeys(key ...string) error {
	err := RDClient.Del(key...).Err()
	return err

}
