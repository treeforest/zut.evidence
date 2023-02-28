package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

const (
	DefaultAddress = "localhost:6379"
	DefaultMaxIdle = 1024
)

func New(address, password string, maxIdle int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   50000,
		IdleTimeout: time.Second * 120,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", address,
				redis.DialConnectTimeout(time.Millisecond*4000),
				redis.DialReadTimeout(time.Millisecond*2000),
				redis.DialWriteTimeout(time.Millisecond*2000),
				redis.DialPassword(password),
			)
			if err != nil {
				panic(err)
			}
			return conn, nil
		},
	}
}
