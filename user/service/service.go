package service

import (
	"context"
	"github.com/jayvib/golog"
	"golang.org/x/crypto/bcrypt"
	"gophr.v2/user"
	"gophr.v2/gophr.client/gophr.api"
)

func New(repo user.Repository) *Service {
	return &Service{repo: repo}
}

type Service struct {
	repo user.Repository
}

func (s *Service) GetByID(ctx context.Context, id string) (*user.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) GetByUsername(ctx context.Context, uname string) (*user.User, error) {
	return s.repo.GetByUsername(ctx, uname)
}

func (s *Service) Save(ctx context.Context, usr *user.User) error {
	return s.repo.Save(ctx, usr)
}

func (s *Service) GetAndComparePassword(ctx context.Context, username, password string) (*user.User, error) {
	// Get the users information
	usr, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// compare the password
	golog.Debug(username)
	golog.Debug(password)
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		golog.Debug(err)
		// if not match then return ErrorCredential
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, gophr.ErrorInvalidCredentials
		}
		return nil, err
	}

	// if match then return the user's information excluding the password
	usr.Password = ""

	return usr, nil
}
