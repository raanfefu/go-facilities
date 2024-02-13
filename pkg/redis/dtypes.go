package redis

import "context"

type RedisClient interface {
	Context() context.Context
	Init()
	Get(ctx context.Context, key string) (string, error)
	PingResult() (string, error)
}
