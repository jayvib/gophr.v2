//+build integration

package file_test

import (
  "context"
  "encoding/json"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "gophr.v2/session"
  "gophr.v2/session/repository/file"
  "os"
  "testing"
  "time"
)

func TestRepository_Save(t *testing.T) {
  sess := &session.Session{
    ID: "sesstest",
    UserID: "userid",
    Expiry: time.Now(),
  }

  filename := "./session_test.db"
  repo := file.New(filename)
  err := repo.Save(context.Background(), sess)
  assert.NoError(t, err)
  assertSessionFound(err, t, filename,sess)
}

func TestRepository_Delete(t *testing.T) {
  sess := &session.Session{
    ID: "sesstest",
    UserID: "userid",
    Expiry: time.Now(),
  }

  filename := "./session_test.db"
  repo := file.New(filename)
  saveToSession(t, repo, sess)
  err := repo.Delete(context.Background(), sess.ID)
  assert.NoError(t, err)
  assertSessionNotFound(err, t, filename,sess)
}

func TestRepository_Find(t *testing.T) {
  sess := &session.Session{
    ID: "sesstest",
    UserID: "userid",
    Expiry: time.Now(),
  }

  filename := "./session_test.db"
  repo := file.New(filename)
  saveToSession(t, repo, sess)

  got, err := repo.Find(context.Background(), sess.ID)
  assert.NoError(t, err)
  assert.Equal(t, sess, got)
  teardownFile(t, filename)
}

func saveToSession(t *testing.T, repo session.Repository, sess *session.Session) {
  t.Helper()
  err := repo.Save(context.Background(), sess)
  assert.NoError(t, err)
}

func assertSessionFound(err error, t *testing.T, filename string, sess *session.Session) {
  t.Helper()
  file, err := os.Open(filename)
  assert.NoError(t, err)
  var sessions map[string]*session.Session
  err = json.NewDecoder(file).Decode(&sessions)
  assert.NoError(t, err)
  _, ok := sessions[sess.ID]
  assert.True(t, ok)
  err = file.Close()
  require.NoError(t, err)
  teardownFile(t, file.Name())
}

func teardownFile(t *testing.T, filename string) {
  err := os.Remove(filename)
  assert.NoError(t, err)
}



func assertSessionNotFound(err error, t *testing.T, filename string, sess *session.Session) {
  t.Helper()
  file, err := os.Open(filename)
  assert.NoError(t, err)
  var sessions map[string]*session.Session
  err = json.NewDecoder(file).Decode(&sessions)
  assert.NoError(t, err)
  _, ok := sessions[sess.ID]
  assert.False(t, ok)
  err = file.Close()
  require.NoError(t, err)
  teardownFile(t , file.Name())
}
