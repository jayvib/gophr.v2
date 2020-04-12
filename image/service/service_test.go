package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gophr.v2/image"
	"gophr.v2/image/imageutil"
	"gophr.v2/image/mocks"
	"gophr.v2/user/userutil"
	"gophr.v2/util/valueutil"
	"testing"
	"time"
)

var dummyContext = context.Background()

func TestService_Find(t *testing.T) {
	t.Run("Image Found", func(t *testing.T){
		want := &image.Image{
			ID: 1,
			UserID: userutil.GenerateID(),
			ImageID: imageutil.GenerateID(),
			CreatedAt: valueutil.TimePointer(time.Now()),
			Name: "Luffy Monkey",
			Location: "East Blue",
			Size: 1024,
			Description: "A Pirate King from East Blue",
		}

		repo := new(mocks.Repository)
		repo.On("Find", mock.Anything, mock.AnythingOfType("string")).Return(want, nil).Once()

		svc := New(repo)
		got, err := svc.Find(dummyContext, want.ImageID)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Image Not Found", func(t *testing.T){
		repo := new(mocks.Repository)
		repo.On("Find", mock.Anything, mock.AnythingOfType("string")).Return(nil, image.ErrNotFound).Once()

		svc := New(repo)
		_, err := svc.Find(dummyContext, "notexists")
		assert.Error(t, err)
		assert.Equal(t, image.ErrNotFound, err)
	})
}
