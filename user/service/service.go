package service

import (
  "context"
  "github.com/jayvib/golog"
  "golang.org/x/crypto/bcrypt"
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
	usr, err := s.repo.GetByID(ctx, id)
	if err != nil {
	  return nil, user.NewError(err).AddContext("ID", id)
  }
	return usr, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	usr, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
	  return nil, user.NewError(err).AddContext("Email", email)
  }
	return usr, nil
}

func (s *Service) GetByUsername(ctx context.Context, uname string) (*user.User, error) {
	usr, err := s.repo.GetByUsername(ctx, uname)
	if err != nil {
	  return nil, user.NewError(err).AddContext("Username", uname)
  }
	return usr, nil
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
			return nil, user.ErrInvalidCredentials
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

func (s *Service) Update(ctx context.Context, usr *user.User) error {

  // Check first if exists
  _, err := s.repo.GetByID(ctx, usr.ID)
  if err != nil {
    if err == user.ErrNotFound {
      err = user.ErrUserNotExists
    }
    return user.NewError(err).AddContext("ID", usr.ID)
  }

	usr.UpdatedAt = valueutil.TimePointer(time.Now().UTC())
	return s.repo.Update(ctx, usr)
}

func (s *Service) Register(ctx context.Context, usr *user.User) error {
	if err := validateUser(usr); err != nil {
		return user.NewError(err)
	}

  // Check first the usr if already exists
  _, err := s.repo.GetByEmail(ctx, usr.Email)
  if err == nil {
    return user.NewError(user.ErrUserExists)
  }

	usr.CreatedAt = valueutil.TimePointer(time.Now().UTC())
	usr.UserID = userutil.GenerateID()
	// Create a password
	hash, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.NewError(err)
	}

	usr.Password = string(hash)

	return s.repo.Save(ctx, usr)
}

func (s *Service) Login(ctx context.Context, user *user.User) error {
  // TODO: To be implemented
	return nil
}

func validateUser(usr *user.User) error {
	if usr.Username == "" {
		return user.ErrEmptyUsername
	}

	if usr.Email == "" {
		return user.ErrEmptyEmail
	}

	if usr.Password == "" {
	  return user.ErrEmptyPassword
  }
	return nil
}
