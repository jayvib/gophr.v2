package redis

import (
  "context"
  "encoding/json"
  "github.com/go-redis/redis/v8"
  "gophr.v2/session"
  "strings"
)

func New(client *redis.Client) session.Repository {
  return &repository{client: client}
}

type repository struct {
  client *redis.Client
}

func (r *repository) Find(ctx context.Context, id string) (*session.Session, error) {
  val, err := r.client.Get(ctx, id).Result()
  if err != nil {
    if strings.Contains(err.Error(), "nil") {
      return nil, session.ErrNotFound
    }
    return nil, err
  }

  strReader := strings.NewReader(val)

  sess := new(session.Session)
  err = json.NewDecoder(strReader).Decode(&sess)
  if err != nil {
    return nil, err
  }

  return sess, nil
}

func  (r *repository) Save(ctx context.Context, s *session.Session) error {
  payload, err := json.Marshal(s)
  if err != nil {
    return err
  }
  _, err = r.client.Set(ctx, s.ID, payload, session.DefaultExpiry).Result()
  if err != nil {
    return err
  }

  return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
  r.client.Del(ctx, id)
  return nil
}