//+build unit


package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gophr.v2/errors"
	"gophr.v2/user"
	"gophr.v2/user/mocks"
)

func TestService_GetByID(t *testing.T) {
	t.Run("Existing user id should return the user information", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := &user.User{
			ID:       12345,
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
		}
		repo.On("GetByID", mock.Anything, mock.AnythingOfType("int")).Return(want, nil)
		svc := New(repo)
		got, _ := svc.GetByID(context.Background(), 12345)
		assert.Equal(t, want, got)
	})

	t.Run("Not existing user should return a ErrNotFound error", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := errors.ErrorNotFound
		repo.On("GetByID", mock.Anything, mock.AnythingOfType("int")).Return(nil, errors.ErrorNotFound)
		svc := New(repo)
		_, got := svc.GetByID(context.Background(), 12345)
		assert.Equal(t, want, got)
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
		want := errors.ErrorNotFound
		repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.ErrorNotFound)
		svc := New(repo)
		_, got := svc.GetByEmail(context.Background(), "12345")
		assert.Equal(t, want, got)
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
		want := errors.ErrorNotFound
		repo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.ErrorNotFound)
		svc := New(repo)
		_, got := svc.GetByUsername(context.Background(), "luffy.monkey")
		assert.Equal(t, want, got)
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
  t.Run("Register User When Not Yet Exists", func(t *testing.T){
    repo := new(mocks.Repository)
    repo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
    repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, ErrNotFound).Once()
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

  t.Run("During Registration Username is Empty", func(t *testing.T){
   repo := new(mocks.Repository)
   svc := New(repo)
   want := &user.User{
     Email:    "luffy.monkey@gmail.com",
     Password: "iampirateking",
   }

   input := *want
   err := svc.Register(context.Background(), &input)
   assert.Error(t, err)
   assert.Equal(t, ErrUsernameEmpty, err)
  })

  t.Run("During Registration Email is Empty", func(t *testing.T){
    repo := new(mocks.Repository)
    svc := New(repo)
    want := &user.User{
      Username: "luffy.monkey",
      Password: "iampirateking",
    }

    input := *want
    err := svc.Register(context.Background(), &input)
    assert.Error(t, err)
    assert.Equal(t, ErrEmptyEmail, err)
  })

  t.Run("During Registration Password is Empty", func(t *testing.T){
    repo := new(mocks.Repository)
    svc := New(repo)
    want := &user.User{
      Email: "luffy.monkey@gmail.com",
      Username: "luffy.monkey",
    }

    input := *want
    err := svc.Register(context.Background(), &input)
    assert.Error(t, err)
    assert.Equal(t, ErrEmptyPassword, err)
  })

  t.Run("Register User When Already Exists", func(t *testing.T){
    res := &user.User{
      Username: "luffy.monkey",
      Email:    "luffy.monkey@gmail.com",
      Password: "iampirateking",
    }

    repo := new(mocks.Repository)
    repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(res, nil).Once()
    svc := New(repo)

    err := svc.Register(context.Background(), res)
    assert.Error(t, ErrUserExists, err)
    repo.AssertExpectations(t)
  })
}

func TestService_Login(t *testing.T) {

}


func TestService_Update(t *testing.T) {
  repo := new(mocks.Repository)
  repo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
  svc := New(repo)

  want := &user.User{
    ID:       12345,
    Username: "luffy.monkey",
    Email:    "luffy.monkey@gmail.com",
  }

  err := svc.Update(context.Background(), want)
  assert.NotEmpty(t, want.UpdatedAt)
  assert.NoError(t, err)
  repo.AssertExpectations(t)
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
  // TODO: Please Implement Me
}

func TestService_GetAndComparePassword(t *testing.T) {
  // TODO: Please Implement me
}
