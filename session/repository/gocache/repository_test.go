package gocache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophr.v2/session"
	"gophr.v2/session/sessionutil"
	"gophr.v2/util/randutil"
	"testing"
	"time"
)

var defaultCtx = context.Background()

func TestRepository_Find(t *testing.T) {
	c := cache.New(defaultExpirationTime, 10*time.Minute)
	t.Run("Found", func(t *testing.T){
		want := &session.Session{
			ID: sessionutil.GenerateID(),
			UserID: randutil.GenerateID("user"),
			Expiry: time.Now(),
		}
		err := c.Add(want.ID, want, defaultExpirationTime)
		require.NoError(t, err)

		r := New(c)
		got, _ := r.Find(defaultCtx, want.ID)
		assert.Equal(t, want, got)
	})

	t.Run("Not Found", func(t *testing.T){
		r := New(c)
		got, err := r.Find(defaultCtx, "notexists")
		assert.Nil(t, got)
		assert.Error(t, err)
	})
}

func TestRepository_Save(t *testing.T) {
}

func TestRepository_Delete(t *testing.T) {
}