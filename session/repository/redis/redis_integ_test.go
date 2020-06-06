// +build integration

package redis_test

import (
	"context"
	"encoding/json"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/jayvib/golog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophr.v2/config"
	"gophr.v2/driver/redis"
	"gophr.v2/session"
	sessionrepo "gophr.v2/session/repository/redis"
	"gophr.v2/session/sessionutil"
	"gophr.v2/user/userutil"
	"testing"
)

var conf = &config.Config{
	Redis: config.Redis{
		Address:  "localhost:6379",
		Username: "", Password: "",
		Database: 0,
	},
}

var dummyCtx = context.Background()

func TestRepository_Find(t *testing.T) {
	golog.Debug(golog.DebugLevel)
	client := redis.New(conf)

	repo := sessionrepo.New(client)
	t.Run("Found", func(t *testing.T) {
		// Set data
		want := &session.Session{
			ID:     sessionutil.GenerateID(),
			UserID: userutil.GenerateID(),
		}
		preSaveToRedis(t, client, want)
		got, err := repo.Find(dummyCtx, want.ID)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Not Found", func(t *testing.T) {
		_, err := repo.Find(dummyCtx, "notexists")
		assert.Error(t, err)
		assert.Equal(t, session.ErrNotFound, err)
	})

}

func TestRepository_Save(t *testing.T) {
	client := redis.New(conf)
	want := &session.Session{
		ID:     sessionutil.GenerateID(),
		UserID: userutil.GenerateID(),
	}

	repo := sessionrepo.New(client)
	err := repo.Save(dummyCtx, want)
	require.NoError(t, err)

	got, _ := repo.Find(dummyCtx, want.ID)
	assert.Equal(t, want, got)
}

func preSaveToRedis(t *testing.T, client *redis2.Client, want *session.Session) {
	t.Helper()
	payload, err := json.Marshal(want)
	require.NoError(t, err)
	_, err = client.Set(dummyCtx, want.ID, payload, session.DefaultExpiry).Result()
	require.NoError(t, err)
}
