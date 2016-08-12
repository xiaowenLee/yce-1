package redigo

import (
	redis "github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

const (
	REDIS_SERVER = "172.21.1.11:32379"

	MAX_IDLE   = 1
	MAX_ACTIVE = 10

	IDEL_TIMEOUT = 180 * time.Second
)

var once sync.Once

var instance *redis.Pool

func NewRedisClient() *redis.Pool {
	once.Do(func() {
		instance = &redis.Pool{
			MaxIdle:     MAX_IDLE,
			MaxActive:   MAX_ACTIVE,
			IdleTimeout: IDEL_TIMEOUT,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", REDIS_SERVER)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		} // instance
	})

	return instance
}
