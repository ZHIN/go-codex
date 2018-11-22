package rds

import (
	"github.com/go-redis/redis"
)

var Default *redis.Client

func SetDefaultClient(option *redis.Options) {
	Default = redis.NewClient(option)
}
