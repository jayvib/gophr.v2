// +build integration

package redis_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gophr.v2/config"
	"gophr.v2/driver/redis"
	"testing"
)

func TestNew(t *testing.T) {
	client := redis.New(&config.Config{
		Redis: config.Redis{
			Address:  "localhost:6379",
			Username: "", Password: "",
			Database: 0,
		},
	})
	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	assert.NoError(t, err)
	assert.Equal(t, "PONG", pong)
}

func ExampleNewClient() {
	client := redis.New(&config.Config{
		Redis: config.Redis{
			Address:  "localhost:6379",
			Username: "", Password: "",
			Database: 0,
		},
	})
	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	fmt.Println(pong, err)
}
