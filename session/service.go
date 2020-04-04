package session

import "context"

//go:generate mockery --name=Service

type Service interface {
  Find(ctx context.Context, id string) (*Session, error)
  Save(ctx context.Context, s *Session) error
  Delete(ctx context.Context, s *Session) error
}
