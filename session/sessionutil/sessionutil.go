package sessionutil

import (
  "gophr.v2/session"
  "gophr.v2/util/randutil"
  "net/http"
  "time"
)

// GenerateID generates a random tokens based on UUID Version 5
func GenerateID() string {
  return randutil.GenerateID("session")
}

func WriteSessionTo(w http.ResponseWriter) *session.Session {
  expiry := time.Now().Add(session.Duration)
  sess := &session.Session{
    ID: GenerateID(),
    Expiry: expiry,
  }

  cookie := http.Cookie{
    Name: session.CookieName,
    Value: sess.ID,
    Expires: sess.Expiry,
  }

  http.SetCookie(w, &cookie)

  return sess
}
