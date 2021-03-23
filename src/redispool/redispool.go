package redispool

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"time"
)

var pool *redigo.Pool

func init() {
	redisHost := "127.0.0.1"
	redisPort := 6379
	poolSize := 20
	pool = &redigo.Pool{
		MaxIdle:     poolSize,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", fmt.Sprintf("%s:%d", redisHost, redisPort))
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

func Get() redigo.Conn {
	return pool.Get()
}
