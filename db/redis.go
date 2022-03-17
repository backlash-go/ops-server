package db

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

var dataPool *redis.Pool

func GetRedisPool() *redis.Pool {
	return pool
}

func GetDataRedisPool() *redis.Pool {
	return dataPool
}

func InitRedis(addr, pwd, choiceDB string) {
	pool = &redis.Pool{
		MaxIdle:     10000,
		MaxActive:   10000,
		IdleTimeout: 5 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if pwd == "" {
				return c, nil
			}
			if _, err := c.Do("AUTH", pwd); err != nil {
				_ = c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", choiceDB); err != nil {
				_ = c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}

func InitDataRedis(addr, pwd, choiceDB string) {
	dataPool = &redis.Pool{
		MaxIdle:     10000,
		MaxActive:   10000,
		IdleTimeout: 5 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if pwd == "" {
				return c, nil
			}
			if _, err := c.Do("AUTH", pwd); err != nil {
				_ = c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", choiceDB); err != nil {
				_ = c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}
