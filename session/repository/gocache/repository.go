package gocache

import (
	"context"
	"fmt"
	"gophr.v2/session"
	"time"
)

const defaultExpirationTime = 24*time.Hour

type CacheIFace interface {
	Get(id string) (data interface{}, ok bool)
	Add(id string, data interface{}, d time.Duration) error
	Delete(id string)
}

func New(c CacheIFace) *Repository {
	return &Repository{
		c: c,
	}
}

type Repository struct {
	c CacheIFace
}

type result struct {
	sess interface{}
	err error
}

func (r *Repository) Find(ctx context.Context, id string) (*session.Session, error) {
	res := make(chan result, 1)

	go func() {
		defer close(res)
		var dataRes result

		select {
		case <-ctx.Done():
			dataRes.err = ctx.Err()
		default:
			v, ok := r.c.Get(id)
			if !ok {
				dataRes.err = session.ErrNotFound
			}
			dataRes.sess = v
		}
		res <- dataRes
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-res:
		if res.err != nil {
			return nil, res.err
		}
		return res.sess.(*session.Session), nil
	}
}

func (r *Repository) Save(ctx context.Context, s *session.Session) error {
	res := make(chan error, 1)
	go func() {
		defer close(res)
		var err error

		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
			e := r.c.Add(s.ID, s, defaultExpirationTime)
			if e != nil {
				err = fmt.Errorf("%w:%s", session.ErrItemExists, e.Error())
			}
		}
		res <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-res:
		return err
	}
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	return nil
}