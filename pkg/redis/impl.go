package redis

import (
	"context"
	"fmt"
	"log"

	v9 "github.com/redis/go-redis/v9"
)

func (c impl) Context() context.Context {
	return context.Background()
}

func (r *impl) Init() {
	options := &v9.Options{
		Addr: fmt.Sprintf("%s:%d", r.Parameters.Host, r.Parameters.Port),
	}
	if r.Parameters.RequiredPass {
		options.Username = r.Parameters.Username
		options.Password = r.Parameters.Password
	}

	r.client = v9.NewClient(options)

	_, err := r.PingResult()
	if err != nil {
		log.Printf("Error ping redis: %s", err)
	}
	log.Println("Initializing redis server... Done ✓")
}

func (r *impl) Status() {
}

/*
	func (r *impl) Del(ctx context.Context, keys ...string) error {
		return nil
	}
*/
func (r *impl) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(r.Context(), key).Result()
}

func (r *impl) PingResult() (string, error) {
	status := r.client.Ping(r.Context())
	return status.String(), status.Err()
}

func (r *impl) SetParametersConfiguration(parameters Parameters) {
	r.Parameters = parameters
}
