//+build unit

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

		svc := New(repo, nil)
		got, err := svc.Find(dummyContext, want.ImageID)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Image Not Found", func(t *testing.T){
		repo := new(mocks.Repository)
		repo.On("Find", mock.Anything, mock.AnythingOfType("string")).Return(nil, image.ErrNotFound).Once()

		svc := New(repo, nil)
		_, err := svc.Find(dummyContext, "notexists")
		assert.Error(t, err)
		assert.Equal(t, image.ErrNotFound, err)
	})
}

func TestService_Save(t *testing.T) {
	want := &image.Image{
		UserID: userutil.GenerateID(),
		Name: "Luffy Monkey",
		Location: "East Blue",
		Size: 1024,
		Description: "A Pirate King from East Blue",
	}
	repo := new(mocks.Repository)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*image.Image")).Return(nil).Once()
	svc := New(repo, nil)
	err := svc.Save(dummyContext, want)
	assert.NoError(t, err)
	assert.NotEmpty(t, want.ImageID)
	assert.NotEmpty(t, want.CreatedAt)
	repo.AssertExpectations(t)
}

func TestService_FindAll(t *testing.T) {
	images := []*image.Image{
		{
			CreatedAt: valueutil.TimePointer(time.Now()),
			UserID: userutil.GenerateID(),
			ImageID: imageutil.GenerateID(),
			Name: "Luffy Monkey",
			Location: "East Blue",
			Size: 1024,
			Description: "A Pirate King from East Blue",
		},
		{
			CreatedAt: valueutil.TimePointer(time.Now()),
			UserID: userutil.GenerateID(),
			ImageID: imageutil.GenerateID(),
			Name: "Roronoa Zoro",
			Location: "East Blue",
			Size: 1024,
			Description: "A Swordsman from East Blue",
		},
		{
			CreatedAt: valueutil.TimePointer(time.Now()),
			UserID: userutil.GenerateID(),
			ImageID: imageutil.GenerateID(),
			Name: "Sanji Vinsmoke",
			Location: "West Blue",
			Size: 1024,
			Description: "A Cook from West Blue",
		},
	}
	repo := new(mocks.Repository)
	repo.On("FindAll", mock.Anything, mock.AnythingOfType("int")).Return(images, nil).Once()
	svc := New(repo, nil)
	got, err := svc.FindAll(dummyContext, 0)
	assert.NoError(t, err)
	assert.Len(t, got, 3)
	repo.AssertExpectations(t)
}

func TestService_FindAllByUser(t *testing.T) {
	userId := userutil.GenerateID()
	images := []*image.Image{
		{
			CreatedAt: valueutil.TimePointer(time.Now()),
			UserID: userId,
			ImageID: imageutil.GenerateID(),
			Name: "Luffy Monkey",
			Location: "East Blue",
			Size: 1024,
			Description: "A Pirate King from East Blue",
		},
		{
			CreatedAt: valueutil.TimePointer(time.Now()),
			UserID: userId,
			ImageID: imageutil.GenerateID(),
			Name: "Roronoa Zoro",
			Location: "East Blue",
			Size: 1024,
			Description: "A Swordsman from East Blue",
		},
		{
			CreatedAt: valueutil.TimePointer(time.Now()),
			UserID: userId,
			ImageID: imageutil.GenerateID(),
			Name: "Sanji Vinsmoke",
			Location: "West Blue",
			Size: 1024,
			Description: "A Cook from West Blue",
		},
	}
	repo := new(mocks.Repository)
	repo.On("FindAllByUser", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int")).Return(images, nil).Once()
	svc := New(repo, nil)
	got, err := svc.FindAllByUser(dummyContext, userId, 0)
	assert.NoError(t, err)
	assert.Len(t, got, 3)
	repo.AssertExpectations(t)
}