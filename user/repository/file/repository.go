package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jayvib/golog"
	"gophr.v2/user"
	"io"
	"io/ioutil"
	"os"
)

func New(filename string) *FileUserStore {
	file, err := os.Open(filename)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			panic(err)
		}
	}

	s := &FileUserStore{
		filename: filename,
		users:    make(map[string]*user.User),
	}

	// meaning this is a path error not exists
	if err != nil {
		return s
	}

	err = json.NewDecoder(file).Decode(&s.users)
	if err != nil && err != io.EOF {
		panic(err)
	}
	return s
}

type FileUserStore struct {
	filename string
	users    map[string]*user.User
	user.Repository
}

func (s *FileUserStore) GetByID(ctx context.Context, id interface{}) (*user.User, error) {
	usr, ok := s.users[id.(string)]
	if !ok {
		return nil, user.ErrNotFound
	}
	return usr, nil
}
func (s *FileUserStore) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	for _, usr := range s.users {
		if usr.Email == email {
			return usr, nil
		}
	}
	return nil, user.ErrNotFound
}
func (s *FileUserStore) GetByUsername(ctx context.Context, uname string) (*user.User, error) {
	for _, usr := range s.users {
		if usr.Username == uname {
			return usr, nil
		}
	}
	return nil, user.ErrNotFound
}
func (s *FileUserStore) Save(ctx context.Context, usr *user.User) error {
	const op = "FileUserStore.Save"
	// check first if the username is already exists
	res, err := s.GetByUsername(ctx, usr.Username)
	if err == nil {
		golog.Debug(err)
		golog.Debugf("%#v\n", res)
		return user.ErrUserNameExists
	}

	_, err = s.GetByEmail(ctx, usr.Email)
	if err == nil {
		golog.Debug(err)
		return user.ErrEmailExists
	}
	s.users[fmt.Sprintf("%d", usr.ID)] = usr

	content, err := json.MarshalIndent(s.users, "", "	")
	if err != nil {
		return fmt.Errorf("%s: error while marsalling user: %w", op, err)
	}

	err = ioutil.WriteFile(s.filename, content, 0666)
	if err != nil {
		return fmt.Errorf("%s: error while writing to file: %w", op, err)
	}
	return nil
}
