package gocache

import (
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"gophr.v2/session"
	"gophr.v2/session/sessionutil"
	"gophr.v2/util/randutil"
	"testing"
	"time"
)

var defaultCtx = context.Background()

type stubCache struct {
	CacheIFace
}

func (s *stubCache) Get(id string) (interface{}, bool) {
	time.Sleep(500*time.Millisecond)
	return nil, false
}

func (s *stubCache) Add(id string, data interface{}, d time.Duration) error {
	time.Sleep(500*time.Millisecond)
	return nil
}

func TestRepository_Find(t *testing.T) {
	c := cache.New(defaultExpirationTime, 10*time.Minute)
	t.Run("Found", func(t *testing.T){
		want := &session.Session{
			ID: sessionutil.GenerateID(),
			UserID: randutil.GenerateID("user"),
			Expiry: time.Now(),
		}

		r := New(c)

		saveData(t, want, r)

		got, _ := r.Find(defaultCtx, want.ID)
		assert.Equal(t, want, got)
	})

	t.Run("Not Found", func(t *testing.T) {
		r := New(c)
		got, err := r.Find(defaultCtx, "notexists")
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("With Context Cancellation", func(t *testing.T) {
		stub := new(stubCache)
		r := New(stub)
		ctx, cancel := context.WithCancel(defaultCtx)
		time.AfterFunc(10*time.Millisecond, cancel)
		_, err := r.Find(ctx, "notexists")
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})
}

func saveData(t *testing.T, want *session.Session, r *Repository) {
	t.Helper()
	err := r.Save(defaultCtx, want)
	assert.NoError(t, err)
}

func TestRepository_Save(t *testing.T) {
	t.Run("Item Not Exists", func(t *testing.T){
		c := cache.New(defaultExpirationTime, 10*time.Minute)
		want := &session.Session{
			ID: sessionutil.GenerateID(),
			UserID: randutil.GenerateID("user"),
			Expiry: time.Now(),
		}

		r := New(c)
		_ = r.Save(defaultCtx, want)

		// Assert
		assertSavedSession(t, r, want)
	})

	t.Run("Item Already Exists", func(t *testing.T){
		c := cache.New(defaultExpirationTime, 10*time.Minute)
		want := &session.Session{
			ID: sessionutil.GenerateID(),
			UserID: randutil.GenerateID("user"),
			Expiry: time.Now(),
		}

		r := New(c)
		saveData(t, want, r)

		err := r.Save(defaultCtx, want)
		assert.Error(t, err)
		assert.Equal(t, session.ErrItemExists, errors.Unwrap(err))
	})

	t.Run("Cancellation", func(t *testing.T){
		stub := new(stubCache)
		r := New(stub)
		ctx, cancel := context.WithCancel(defaultCtx)
		time.AfterFunc(10*time.Millisecond, cancel)

		want := &session.Session{
			ID: sessionutil.GenerateID(),
			UserID: randutil.GenerateID("user"),
			Expiry: time.Now(),
		}

		err := r.Save(ctx, want)
		assert.Error(t, err)
	})
}

func assertSavedSession(t *testing.T, r *Repository, want *session.Session) {
	t.Helper()
	got, err := r.Find(defaultCtx, want.ID)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}



func TestRepository_Delete(t *testing.T) {
}