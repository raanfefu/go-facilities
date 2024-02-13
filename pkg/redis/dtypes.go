package redis

import "context"

type RedisClient interface {
	Context() context.Context
	Init()
	PingResult() (string, error)
}
