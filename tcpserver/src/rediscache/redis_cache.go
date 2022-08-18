package rediscache

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

var Rds redis.Conn
var Pool *redis.Pool

// RedisPoolInit initialize a redis connectino pool configed by redis_config.go
func RedisPoolInit() *redis.Pool {
	return &redis.Pool{
		MaxIdle: Maxidle,
		MaxActive: Maxactive, 
		IdleTimeout: Idletimeout,
		Wait: true,
		Dial: func() (redis.Conn, error){
			conn, err := redis.Dial("tcp", 
				"0.0.0.0:6379",
				redis.DialReadTimeout(Dialreadtimeout),
				redis.DialWriteTimeout(Dialwritetimeout),
				redis.DialConnectTimeout(Dialconnecttimeout),
			)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			redis.DialDatabase(0)
			return conn, err
		},
	}
}

// init() initialize the connection poll when package is imported
func init() {
	Pool = RedisPoolInit()
}

// RedisInit() return the connection from the connectin pool
func RedisInit() (redis.Conn, error) {
	conn := Pool.Get()
	return conn, nil
}