package file

import (
  "context"
  "encoding/json"
  "gophr.v2/session"
  "io/ioutil"
  "os"
)

func New(filename string) session.Repository {
  r := &repository{
    filename: filename,
    sessions: make(map[string]*session.Session),
  }
  return r
}

type repository struct {
  filename string
  sessions map[string]*session.Session
}

func (r *repository) init() {
  errFunc := func(err error) {
    if err != nil && !os.IsNotExist(err) {
      panic(err)
    }
  }
  file, err := os.Open(r.filename)
  errFunc(err)
  if err == nil {
    defer func() {
      _ = file.Close()
    }()

    err = json.NewDecoder(file).Decode(&r.sessions)
    errFunc(err)
  }
}

func (r *repository) Find(ctx context.Context, id string) (*session.Session, error) {
  var sess *session.Session
  var ok bool
  if sess, ok = r.sessions[id]; !ok {
    return nil, session.ErrNotFound
  }
  return sess, nil
}

func (r *repository) Save(ctx context.Context, s *session.Session) error {
  r.sessions[s.ID] = s
  payload, err := r.marshalToJSON()
  if err != nil {
    return err
  }
  return ioutil.WriteFile(r.filename, payload, 0755)
}

func (r *repository) marshalToJSON() ([]byte, error) {
  payload, err := json.MarshalIndent(r.sessions, "", "  ")
  if err != nil {
    return nil, session.NewError(err, "failed to unmarshal sessions")
  }
  return payload, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
  delete(r.sessions, id)
  payload, err := r.marshalToJSON()
  if err != nil {
    return err
  }
  return ioutil.WriteFile(r.filename, payload, 0755)
}


