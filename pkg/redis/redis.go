package redis

import (
	"github.com/raanfefu/go-facilities/pkg/configure"
	v9 "github.com/redis/go-redis/v9"
)

type impl struct {
	configure.DefaultConfiguraionService
	client     v9.UniversalClient
	Parameters Parameters
}

type Parameters struct {
	Host         string
	Port         uint
	Username     string
	Password     string
	RequiredPass bool
}

func NewRedisClient(conf configure.Configuration) RedisClient {
	obj := &impl{}
	conf.RegistryService("redis", obj)
	return obj
}
