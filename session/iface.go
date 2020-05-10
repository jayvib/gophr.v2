package session

import "context"

// Applying the Interface Segregation Principle

type Finder interface {
	Find(ctx context.Context, id string) (*Session, error)
}

type Saver interface {
	Save(ctx context.Context, s *Session) error
}

type Deleter interface {
	Delete(ctx context.Context, id string) error
}
