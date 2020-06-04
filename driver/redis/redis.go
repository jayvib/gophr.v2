package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jayvib/golog"
	"gophr.v2/config"
)

// New initializes redis client using the passed conf configuration.
// This will panic when the connection is not successful.
func New(conf *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Username: conf.Redis.Username,
		Password: conf.Redis.Password,
		DB:       conf.Redis.Database,
	})

	ctx := context.Background()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	golog.Debugf("REDIS: %s: PING...", pong)

	return client
}
