package service

import (
	"context"
	"github.com/jayvib/golog"
	"golang.org/x/crypto/bcrypt"
	"gophr.v2/errors"
	"gophr.v2/user"
	"gophr.v2/user/userutil"
	"gophr.v2/util/valueutil"
	"time"
)

var _ user.Service = (*Service)(nil)

func New(repo user.Repository) *Service {
	return &Service{repo: repo}
}

type Service struct {
	repo user.Repository
}

func (s *Service) GetByID(ctx context.Context, id interface{}) (*user.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) GetByUsername(ctx context.Context, uname string) (*user.User, error) {
	return s.repo.GetByUsername(ctx, uname)
}

func (s *Service) Save(ctx context.Context, usr *user.User) error {
	usr.CreatedAt = valueutil.TimePointer(time.Now().UTC())
	usr.UserID = userutil.GenerateID()
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
			return nil, errors.ErrorInvalidCredentials
		}
		return nil, err
	}

	// if match then return the user's information excluding the password
	usr.Password = ""

	return usr, nil
}

func (s *Service) GetAll(ctx context.Context, cursor string, num int) (user []*user.User, nextCursor string, err error) {
	return s.repo.GetAll(ctx, cursor, num)
}

func (s *Service) Delete(ctx context.Context, id interface{}) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Update(ctx context.Context, user *user.User) error {
  user.UpdatedAt = valueutil.TimePointer(time.Now().UTC())
	return s.repo.Update(ctx, user)
}

func (s *Service) Register(ctx context.Context, user *user.User) error {
  user.CreatedAt = valueutil.TimePointer(time.Now().UTC())
  user.UserID = userutil.GenerateID()

  // Create a password
 hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
 if err != nil {
   return err
 }

 user.Password = string(hash)

  return s.Save(ctx, user)
}

func (s *Service) Login(ctx context.Context, user *user.User) error {
  return nil
}
