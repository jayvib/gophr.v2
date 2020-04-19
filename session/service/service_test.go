//+build unit

package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gophr.v2/session"
	"gophr.v2/session/mocks"
	"testing"
	"time"
)

func TestService_Delete(t *testing.T) {
	repo := new(mocks.Repository)
	repo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
	svc := New(repo)

	err := svc.Delete(context.Background(), "1234")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestService_Find(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		sess := &session.Session{
			ID:     "test123",
			UserID: "userid123",
			Expiry: time.Now(),
		}

		repo := new(mocks.Repository)
		repo.On("Find", mock.Anything, mock.AnythingOfType("string")).Return(sess, nil).Once()
		svc := New(repo)

		got, _ := svc.Find(context.Background(), sess.ID)
		assert.Equal(t, sess, got)
		repo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		want := session.NewError(session.ErrNotFound, "Failed finding session with ID: 12234")

		repo := new(mocks.Repository)
		repo.On("Find", mock.Anything, mock.AnythingOfType("string")).
			Return(nil, session.ErrNotFound).
			Once()
		svc := New(repo)

		_, err := svc.Find(context.Background(), "12234")
		assert.Error(t, err)
		assert.Equal(t, want, err)
		repo.AssertExpectations(t)
	})
}

func TestService_Save(t *testing.T) {
	sess := &session.Session{
		ID:     "test123",
		UserID: "userid123",
		Expiry: time.Now(),
	}
	repo := new(mocks.Repository)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*session.Session")).Return(nil).Once()
	svc := New(repo)

	err := svc.Save(context.Background(), sess)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
