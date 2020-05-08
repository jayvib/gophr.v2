package gocache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"gophr.v2/session"
	"time"
)

const defaultExpirationTime = 24*time.Hour

func New(c *cache.Cache) *Repository {
	return &Repository{
		c: c,
	}
}

type Repository struct {
	c *cache.Cache
}

func (r *Repository) Find(ctx context.Context, id string) (*session.Session, error) {
	res, ok := r.c.Get(id)
	if !ok {
		return nil, session.ErrNotFound
	}
	return res.(*session.Session), nil
}

func (r *Repository) Save(ctx context.Context, s *session.Session) error {
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	return nil
}