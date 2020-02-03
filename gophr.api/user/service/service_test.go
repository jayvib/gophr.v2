// +build unit

package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gophr.v2/gophr.api/user"
	"gophr.v2/gophr.api/user/mocks"
	"testing"
)

func TestService_GetByID(t *testing.T) {
	t.Run("Existing user id should return the user information", func(t *testing.T){
		repo := new(mocks.Repository)
		want := &user.User{
			ID: "12345",
			Username: "luffy.monkey",
			Email: "luffy.monkey@gmail.com",
		}
		repo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(want, nil)
		svc := New(repo)
		got, _ := svc.GetByID(context.Background(), "12345")
		assert.Equal(t, want, got)
	})

	t.Run("Not existing user should return a ErrNotFound error", func(t *testing.T){
	})
}

func TestService_GetByEmail(t *testing.T) {}

func TestService_GetByUsername(t *testing.T) {}

func TestService_Save(t *testing.T) {}

