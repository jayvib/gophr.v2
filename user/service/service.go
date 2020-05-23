package service

import (
	"context"
	"fmt"
	"github.com/jayvib/golog"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
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

func (s *Service) GetByUserID(ctx context.Context, userId string) (*user.User, error) {
	usr, err := s.repo.GetByUserID(ctx, userId)
	if err != nil {
		return nil, user.NewError(err).AddContext("User ID", userId)
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

func (s *Service) getAndComparePassword(ctx context.Context, username, password string) (*user.User, error) {
	// Get the users information
	usr, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// compare the password
	golog.Debug(username)
	golog.Debug(password)
	golog.Debugf("%#v\n", usr)
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
	_, err := s.repo.GetByUserID(ctx, usr.UserID)
	if err != nil {
		if err == user.ErrNotFound {
			err = user.ErrUserNotExists
		}
		return user.NewError(err).AddContext("ID", usr.UserID)
	}

	usr.UpdatedAt = valueutil.TimePointer(time.Now().UTC())

	// TODO: Hash the password

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

func (s *Service) Login(ctx context.Context, usr *user.User) error {
	// Compare the value of user password and the existing user password
	u, err := s.getAndComparePassword(ctx, usr.Username, usr.Password)
	if err != nil {
		return user.NewError(err)
	}
	usr.Password = ""
	usr.UserID = u.UserID // I don't know if it is right
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

func GetByUserIDs(ctx context.Context, svc user.GetterByUserID, ids ...string) ([]*user.User, error) {
	var users []*user.User

	g, ctx := errgroup.WithContext(ctx)

	type result struct{
		usr *user.User
		err error
		id string
	}

	resultChan := make(chan *result)
	for _, id := range ids {
		id := id
		g.Go(func()error{
			res, err := svc.GetByUserID(ctx, id)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case resultChan <- &result{usr: res, err: err, id: id}:
			}
			return nil
		})
	}

	go func() {
		if err := g.Wait(); err != nil {
			golog.Error(err)
			return
		}
		close(resultChan)
	}()

	errMsg := ""
	for res := range resultChan {
		if res.err == nil {
			users = append(users, res.usr)
		} else {
			if res.err == user.ErrNotFound {
				errMsg += fmt.Sprintf("user with id '%s' not exist\n", res.id)
			} else {
				errMsg += fmt.Sprintf("%s\n", res.err.Error())
			}
		}
	}

	if errMsg != "" {
		return users, errors.New(errMsg)
	}

	return users, nil
}
