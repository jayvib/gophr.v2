package service

import (
  "context"
  "fmt"
  "gophr.v2/session"
)

func New(repo session.Repository) session.Service {
  return &Service{repo: repo}
}

type Service struct {
  repo session.Repository
}

func (s *Service) Find(ctx context.Context, id string) (*session.Session, error) {
  sess, err := s.repo.Find(ctx, id)
  if err != nil {
    return nil, s.wrapError(err, fmt.Sprintf("Failed finding session with ID: %s", id))
  }
  return sess, nil
}

func (s *Service) wrapError(err error, msg string) error {
  if _, ok := err.(*session.Error); !ok {
    err = session.NewError(err, msg)
  }
  return err
}

func (s *Service) Save(ctx context.Context, sess *session.Session) error {
  return s.repo.Save(ctx, sess)
}

func (s *Service) Delete(ctx context.Context, id string) error {
  return s.repo.Delete(ctx, id)
}
