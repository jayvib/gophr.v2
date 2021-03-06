//+build unit

package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"gophr.v2/user/userutil"
	"gophr.v2/util/valueutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gophr.v2/user"
	"gophr.v2/user/mocks"
)

func TestService_GetByID(t *testing.T) {
	t.Run("Existing user id should return the user information", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := &user.User{
			ID:       12345,
			UserID:   userutil.GenerateID(),
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
		}
		repo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(want, nil)
		svc := New(repo)
		got, _ := svc.GetByID(context.Background(), want.ID)
		assert.Equal(t, want, got)
	})

	t.Run("Not existing user should return a ErrNotFound error", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := user.ErrNotFound
		repo.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(nil, user.ErrNotFound)
		svc := New(repo)
		_, got := svc.GetByID(context.Background(), uint(9999))
		require.IsType(t, new(user.Error), got)
		assert.Equal(t, want, errors.Unwrap(got))
	})
}

func TestService_GetByUserID(t *testing.T) {
	t.Run("Existing user id should return the user information", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := &user.User{
			ID:       12345,
			UserID:   userutil.GenerateID(),
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
		}
		repo.On("GetByUserID", mock.Anything, mock.AnythingOfType("string")).Return(want, nil)
		svc := New(repo)
		got, _ := svc.GetByUserID(context.Background(), want.UserID)
		assert.Equal(t, want, got)
	})

	t.Run("Not existing user should return a ErrNotFound error", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := user.ErrNotFound
		repo.On("GetByUserID", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrNotFound)
		svc := New(repo)
		_, got := svc.GetByUserID(context.Background(), "notfoundid")
		require.IsType(t, new(user.Error), got)
		assert.Equal(t, want, errors.Unwrap(got))
	})
}

func TestService_GetByEmail(t *testing.T) {
	t.Run("Existing email should return its equavalent user object", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := &user.User{
			ID:       12345,
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
		}
		repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(want, nil).Once()
		svc := New(repo)
		got, _ := svc.GetByEmail(context.Background(), "luffy.monkey@gmail.com")
		assert.Equal(t, want, got)
	})

	t.Run("Not existing user should return a ErrNotFound error", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := user.ErrNotFound
		repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrNotFound)
		svc := New(repo)
		_, got := svc.GetByEmail(context.Background(), "12345")
		require.IsType(t, new(user.Error), got)
		assert.Equal(t, want, errors.Unwrap(got))
	})
}

func TestService_GetByUsername(t *testing.T) {
	t.Run("Existing username should return its equavalent user object", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := &user.User{
			ID:       12345,
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
		}
		repo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(want, nil).Once()
		svc := New(repo)
		got, _ := svc.GetByUsername(context.Background(), "luffy.monkey")
		assert.Equal(t, want, got)
	})

	t.Run("Not existing username should return a ErrNotFound error", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := user.ErrNotFound
		repo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrNotFound)
		svc := New(repo)
		_, got := svc.GetByUsername(context.Background(), "luffy.monkey")
		require.IsType(t, new(user.Error), got)
		assert.Equal(t, want, errors.Unwrap(got))
	})

}

func TestService_Save(t *testing.T) {
	repo := new(mocks.Repository)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
	svc := New(repo)
	want := &user.User{
		ID:       12345,
		Username: "luffy.monkey",
		Email:    "luffy.monkey@gmail.com",
	}
	err := svc.Save(context.Background(), want)
	assert.NotEmpty(t, want.CreatedAt)
	assert.NotEmpty(t, want.UserID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestService_Register(t *testing.T) {
	t.Run("Register User When Not Yet Exists", func(t *testing.T) {
		repo := new(mocks.Repository)
		repo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
		repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrNotFound).Once()
		svc := New(repo)
		want := &user.User{
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
			Password: "iampirateking",
		}

		input := *want
		err := svc.Register(context.Background(), &input)
		assert.NoError(t, err)
		assert.NotEqual(t, want.Password, input.Password)
		repo.AssertExpectations(t)
	})

	t.Run("During Registration Username is Empty", func(t *testing.T) {
		repo := new(mocks.Repository)
		svc := New(repo)
		want := &user.User{
			Email:    "luffy.monkey@gmail.com",
			Password: "iampirateking",
		}

		input := *want
		err := svc.Register(context.Background(), &input)
		assert.Error(t, err)
		require.IsType(t, new(user.Error), err)
		assert.Equal(t, user.ErrEmptyUsername, errors.Unwrap(err))
		t.Log(err)
	})

	t.Run("During Registration Email is Empty", func(t *testing.T) {
		repo := new(mocks.Repository)
		svc := New(repo)
		want := &user.User{
			Username: "luffy.monkey",
			Password: "iampirateking",
		}

		input := *want
		err := svc.Register(context.Background(), &input)
		assert.Error(t, err)
		assert.IsType(t, new(user.Error), err)
		assert.Equal(t, user.ErrEmptyEmail, errors.Unwrap(err))
	})

	t.Run("During Registration Password is Empty", func(t *testing.T) {
		repo := new(mocks.Repository)
		svc := New(repo)
		want := &user.User{
			Email:    "luffy.monkey@gmail.com",
			Username: "luffy.monkey",
		}

		input := *want
		err := svc.Register(context.Background(), &input)
		assert.Error(t, err)
		assert.IsType(t, new(user.Error), err)
		assert.Equal(t, user.ErrEmptyPassword, errors.Unwrap(err))
	})

	t.Run("Register User When Already Exists", func(t *testing.T) {
		res := &user.User{
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
			Password: "iampirateking",
		}

		repo := new(mocks.Repository)
		repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(res, nil).Once()
		svc := New(repo)

		err := svc.Register(context.Background(), res)
		assert.Error(t, user.ErrUserExists, err)
		repo.AssertExpectations(t)
	})
}

func TestService_Login(t *testing.T) {
	t.Run("Valid Credential", func(t *testing.T) {
		// Simulate registration
		usr := &user.User{
			Email:    "luffy.monkey@gmail.com",
			Username: "luffy.monkey",
			Password: "iampirateking",
		}

		cpy := *usr
		clonedUsr := &cpy
		repo := new(mocks.Repository)
		repo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
		repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrNotFound).Once()
		repo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(clonedUsr, nil).Once()
		svc := New(repo)
		err := svc.Register(context.Background(), clonedUsr)
		require.NoError(t, err)

		// LoginPage
		err = svc.Login(context.Background(), usr)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Invalid Credential", func(t *testing.T) {
		// Simulate registration
		usr := &user.User{
			Email:    "luffy.monkey@gmail.com",
			Username: "luffy.monkey",
			Password: "iampirateking",
		}

		cpy := *usr
		clonedUsr := &cpy

		// Modify to make the password invalid
		usr.Password = "invalidpassword"

		repo := new(mocks.Repository)
		repo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
		repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrNotFound).Once()
		repo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(clonedUsr, nil).Once()
		svc := New(repo)
		err := svc.Register(context.Background(), clonedUsr)
		assert.NoError(t, err)

		// LoginPage
		err = svc.Login(context.Background(), usr)
		require.Error(t, err)

		want := user.NewError(user.ErrInvalidCredentials)
		assert.Equal(t, want, err)
		repo.AssertExpectations(t)
	})

}

func TestService_Update(t *testing.T) {
	t.Run("Updating An Existing User", func(t *testing.T) {
		repo := new(mocks.Repository)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
		repo.On("GetByUserID", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Once()
		svc := New(repo)

		want := &user.User{
			ID:       12345,
			UserID:   userutil.GenerateID(),
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
		}

		err := svc.Update(context.Background(), want)
		assert.NotEmpty(t, want.UpdatedAt)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Updating A Non-Existing User", func(t *testing.T) {
		repo := new(mocks.Repository)
		repo.On("GetByUserID", mock.Anything, mock.AnythingOfType("string")).Return(nil, user.ErrUserNotExists).Once()
		input := &user.User{
			ID:       12345,
			UserID:   userutil.GenerateID(),
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
		}

		svc := New(repo)
		err := svc.Update(context.Background(), input)
		require.Error(t, err)
		assert.IsType(t, new(user.Error), err)
		assert.Equal(t, user.ErrUserNotExists, errors.Unwrap(err))
	})
}

func TestService_Delete(t *testing.T) {
	repo := new(mocks.Repository)
	repo.On("Delete", mock.Anything, mock.AnythingOfType("uint")).Return(nil).Once()
	svc := New(repo)

	want := &user.User{
		ID:       12345,
		Username: "luffy.monkey",
		Email:    "luffy.monkey@gmail.com",
	}

	err := svc.Delete(context.Background(), want.ID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestService_GetAll(t *testing.T) {
	// Add the mock users to the rows
	mockUsers := []*user.User{
		{
			ID:        1,
			UserID:    "testid123",
			Username:  "unit.test",
			Email:     "unit.test@golang.com",
			Password:  "qwerty",
			CreatedAt: valueutil.TimePointer(time.Now()),
			UpdatedAt: valueutil.TimePointer(time.Now()),
		},
		{
			ID:        2,
			UserID:    "testid124",
			Username:  "unit.test01",
			Email:     "unit.test01@golang.com",
			Password:  "qwerty",
			CreatedAt: valueutil.TimePointer(time.Now()),
			UpdatedAt: valueutil.TimePointer(time.Now()),
		},
		{
			ID:        3,
			UserID:    "testid125",
			Username:  "unit.test02",
			Email:     "unit.test02@golang.com",
			Password:  "qwerty",
			CreatedAt: valueutil.TimePointer(time.Now()),
			UpdatedAt: valueutil.TimePointer(time.Now()),
		},
	}

	repo := new(mocks.Repository)
	repo.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int")).Return(mockUsers, "", nil).Once()
	svc := New(repo)

	cursor := userutil.EncodeCursor(*mockUsers[0].CreatedAt)
	got, _, _ := svc.GetAll(context.Background(), cursor, 3)
	assert.Len(t, got, 3)
	repo.AssertExpectations(t)
}

type svcStub struct {
	user.Service
	data map[string]*user.User
}

func (s *svcStub) GetByUserID(ctx context.Context, id string) (*user.User, error) {
	u := s.data[id]
	return u, nil
}

func TestGetByUserIDs(t *testing.T) {
	want := []*user.User{
		{
			ID:        1,
			UserID:    userutil.GenerateID(),
			Username:  "luffy.monkey",
			Email:     "luffy.monkey@gmail.com",
			Password:  "qwqewrt",
			CreatedAt: valueutil.TimePointer(time.Now()),
		},
		{
			ID:        2,
			UserID:    userutil.GenerateID(),
			Username:  "sanji.vinsmoke",
			Email:     "sanji.vinsmoke@gmail.com",
			Password:  "qwqewrt",
			CreatedAt: valueutil.TimePointer(time.Now()),
		},
	}
	svc := &svcStub{
		data: map[string]*user.User{
			want[0].UserID: want[0],
			want[1].UserID: want[1],
		},
	}

	got, err := GetByUserIDs(context.Background(), svc, want[0].UserID, want[1].UserID)
	assert.NoError(t, err)

	assert.Equal(t, want, got)

}
