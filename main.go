package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func main() {

	pool := &redis.Pool{
		// 連線的 callback 定義
		Dial: func() (redis.Conn, error) {

			//建構一條連線
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}

			//在這邊可以做連線池初始化 選擇 redis db的動作
			// if _, err := c.Do("SELECT", db); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			return c, nil
		},

		//定期對 redis server 做 ping/pong 測試
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	//這邊要非常注意，redis用完連線，請一定要做close的動作，否則有機會造成 memory leak

	// close完的連線，會回到 connection pool
	conn := pool.Get()
	defer conn.Close()

	// 寫入 redis
	reply, err := conn.Do("MSET", "key1", "hello", "key2", "world")
	if err != nil {

	}

	fmt.Println("reply =>", reply)

	// 讀取 redis
	response, err := redis.Strings(conn.Do("MGET", "key1", "key2"))
	if err != nil {

	}

	fmt.Println("response =>", response)
}
