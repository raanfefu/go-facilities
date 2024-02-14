package redis

import "context"

type RedisClient interface {
	Context() context.Context
	SetParametersConfiguration(params Parameters)
	Init()
	Get(ctx context.Context, key string) (string, error)
	PingResult() (string, error)
}
