package rediscache

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

var Rds redis.Conn
var Pool *redis.Pool
// type User struct{
// 	Name string `json:"name"`
// 	Data struct{
// 		Password string `json:"password"`
// 	}	`json:"data"`
// }

func RedisPoolInit() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 5,
		MaxActive: 0, 
		IdleTimeout: 1 * time.Second,
		Dial: func() (redis.Conn, error){
			conn, err := redis.Dial("tcp", "0.0.0.0:6379")
			if err != nil {
				log.Println(err)
				return nil, err
			}
			redis.DialDatabase(0)
			return conn, err
		},
	}
}

func init() {
	Pool = RedisPoolInit()
}

func RedisInit() redis.Conn {
	return Pool.Get()
}


// func main() {
// 	RedisInit()
// 	defer RedisClose()

// 	user := User{
// 		Name: "ding",
// 		Data: struct{Password string "json:\"password\""}{Password: "123"},
// 	}

// 	user_encode, err := json.Marshal(user)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	_, err = rds.Do("set", "user1", user_encode)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	user_get, err := redis.Bytes(rds.Do("get", "user1"))
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	user_return := new(User)
// 	err = json.Unmarshal(user_get, user_return)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	fmt.Printf("return value: \n %+v \n", *user_return)
// }