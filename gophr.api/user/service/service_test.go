// +build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gophr.v2/gophr.api/errors"
	"gophr.v2/gophr.api/user"
	"gophr.v2/gophr.api/user/mocks"
)

func TestService_GetByID(t *testing.T) {
	t.Run("Existing user id should return the user information", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := &user.User{
			ID:       "12345",
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
		}
		repo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(want, nil)
		svc := New(repo)
		got, _ := svc.GetByID(context.Background(), "12345")
		assert.Equal(t, want, got)
	})

	t.Run("Not existing user should return a ErrNotFound error", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := errors.ErrorNotFound
		repo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.ErrorNotFound)
		svc := New(repo)
		_, got := svc.GetByID(context.Background(), "12345")
		assert.Equal(t, want, got)
	})
}

func TestService_GetByEmail(t *testing.T) {
	t.Run("Existing email should return its equavalent user object", func(t *testing.T) {
		repo := new(mocks.Repository)
		want := &user.User{
			ID:       "12345",
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
		}
		repo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(want, nil).Once()
		svc := New(repo)
		got, _ := svc.GetByEmail(context.Background(), "12345")
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
	t.Run("Existing username should return its equavalent user object", func(t *testing.T){
		repo := new(mocks.Repository)
		want := &user.User{
			ID:       "12345",
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
}
