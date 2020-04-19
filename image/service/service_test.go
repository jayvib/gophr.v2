//+build unit

package service

import (
	"bufio"
	"bytes"
	"context"
	"github.com/jayvib/golog"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gophr.v2/http/httputil"
	"gophr.v2/image"
	"gophr.v2/image/imageutil"
	"gophr.v2/image/mocks"
	"gophr.v2/user/userutil"
	"gophr.v2/util/valueutil"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var dummyContext = context.Background()

func TestService_Find(t *testing.T) {
	t.Run("Image Found", func(t *testing.T) {
		want := &image.Image{
			ID:          1,
			UserID:      userutil.GenerateID(),
			ImageID:     imageutil.GenerateID(),
			CreatedAt:   valueutil.TimePointer(time.Now()),
			Name:        "Luffy Monkey",
			Location:    "East Blue",
			Size:        1024,
			Description: "A Pirate King from East Blue",
		}

		repo := new(mocks.Repository)
		repo.On("Find", mock.Anything, mock.AnythingOfType("string")).Return(want, nil).Once()

		svc := New(repo, nil, nil)
		got, err := svc.Find(dummyContext, want.ImageID)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Image Not Found", func(t *testing.T) {
		repo := new(mocks.Repository)
		repo.On("Find", mock.Anything, mock.AnythingOfType("string")).Return(nil, image.ErrNotFound).Once()

		svc := New(repo, nil, nil)
		_, err := svc.Find(dummyContext, "notexists")
		assert.Error(t, err)
		assert.Equal(t, image.ErrNotFound, err)
	})
}

func TestService_Save(t *testing.T) {
	want := &image.Image{
		UserID:      userutil.GenerateID(),
		Name:        "Luffy Monkey",
		Location:    "East Blue",
		Size:        1024,
		Description: "A Pirate King from East Blue",
	}
	repo := new(mocks.Repository)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*image.Image")).Return(nil).Once()
	svc := New(repo, nil, nil)
	err := svc.Save(dummyContext, want)
	assert.NoError(t, err)
	assert.NotEmpty(t, want.ImageID)
	assert.NotEmpty(t, want.CreatedAt)
	repo.AssertExpectations(t)
}

func TestService_FindAll(t *testing.T) {
	images := []*image.Image{
		{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userutil.GenerateID(),
			ImageID:     imageutil.GenerateID(),
			Name:        "Luffy Monkey",
			Location:    "East Blue",
			Size:        1024,
			Description: "A Pirate King from East Blue",
		},
		{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userutil.GenerateID(),
			ImageID:     imageutil.GenerateID(),
			Name:        "Roronoa Zoro",
			Location:    "East Blue",
			Size:        1024,
			Description: "A Swordsman from East Blue",
		},
		{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userutil.GenerateID(),
			ImageID:     imageutil.GenerateID(),
			Name:        "Sanji Vinsmoke",
			Location:    "West Blue",
			Size:        1024,
			Description: "A Cook from West Blue",
		},
	}
	repo := new(mocks.Repository)
	repo.On("FindAll", mock.Anything, mock.AnythingOfType("int")).Return(images, nil).Once()
	svc := New(repo, nil, nil)
	got, err := svc.FindAll(dummyContext, 0)
	assert.NoError(t, err)
	assert.Len(t, got, 3)
	repo.AssertExpectations(t)
}

func TestService_FindAllByUser(t *testing.T) {
	userId := userutil.GenerateID()
	images := []*image.Image{
		{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userId,
			ImageID:     imageutil.GenerateID(),
			Name:        "Luffy Monkey",
			Location:    "East Blue",
			Size:        1024,
			Description: "A Pirate King from East Blue",
		},
		{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userId,
			ImageID:     imageutil.GenerateID(),
			Name:        "Roronoa Zoro",
			Location:    "East Blue",
			Size:        1024,
			Description: "A Swordsman from East Blue",
		},
		{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userId,
			ImageID:     imageutil.GenerateID(),
			Name:        "Sanji Vinsmoke",
			Location:    "West Blue",
			Size:        1024,
			Description: "A Cook from West Blue",
		},
	}
	repo := new(mocks.Repository)
	repo.On("FindAllByUser", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int")).Return(images, nil).Once()
	svc := New(repo, nil, nil)
	got, err := svc.FindAllByUser(dummyContext, userId, 0)
	assert.NoError(t, err)
	assert.Len(t, got, 3)
	repo.AssertExpectations(t)
}

func TestService_CreateImageFromURL(t *testing.T) {
	golog.SetLevel(golog.DebugLevel)

	// Learned remove server stubbing here:
	// https://itnext.io/how-to-stub-requests-to-remote-hosts-with-go-6c2c1db32bf2
	dummyImage, client, teardown := setupServerAndClient(t)
	defer teardown()

	t.Run("Fetching the Image From Remote", func(t *testing.T) {
		repo := new(mocks.Repository)
		repo.On("Save", mock.Anything, mock.AnythingOfType("*image.Image")).Return(nil).Once()
		dummyFs := afero.NewMemMapFs()
		svc := New(repo, dummyFs, client)
		dummyUserID := "qwerty1234"
		dummyDescription := "A Unit Testing"
		got, err := svc.CreateImageFromURL(dummyContext, "http://127.0.0.1/image.png", dummyUserID, dummyDescription)
		assert.NoError(t, err)
		assertImage(t, got, int64(len(dummyImage)), "image.png", ".png", dummyUserID, dummyDescription)
		assertImageContent(err, dummyFs, got, t, dummyImage)
		repo.AssertExpectations(t)
	})

	t.Run("Error while Doing Client Request", func(t *testing.T) {
		url := "127.0.0.1:12345/testingerror"
		svc := New(nil, nil, client)
		_, err := svc.CreateImageFromURL(dummyContext, url, "user123", "testing the client request error")
		assert.Error(t, err)
		assert.Equal(t, image.ErrInvalidImageURL, err)
	})

	t.Run("Failed Request To Remote", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/image.png", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		client, teardown = httputil.DummyClient(mux)
		defer teardown()

		url := "http://127.0.0.1/image.png"
		svc := New(nil, nil, client)
		_, err := svc.CreateImageFromURL(dummyContext, url, "user123", "testing the client request failed status")
		assert.Error(t, err)
		assert.Equal(t, image.ErrFailedRequest, err)
	})

	t.Run("Invalid Content Type", func(t *testing.T) {
		// TODO: Pls implement
	})
}

func TestService_CreateImageFromFile(t *testing.T) {
	stat, err := os.Stat("./testdata/simple.png")
	require.NoError(t, err)
	f, err := os.Open("./testdata/simple.png")
	require.NoError(t, err)
	defer f.Close()

	dummyFs := afero.NewMemMapFs()
	repo := new(mocks.Repository)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*image.Image")).Return(nil).Once()
	svc := New(repo, dummyFs, nil)
	got, err := svc.CreateImageFromFile(dummyContext, f, "simple.png", "A Unit Test", "user12345")
	assert.NoError(t, err)
	assertImage(t, got, stat.Size(), "simple.png", ".png", "user12345", "A Unit Test")
	repo.AssertExpectations(t)
}

func setupServerAndClient(t *testing.T) ([]byte, *http.Client, func()) {
	dummyImage, err := ioutil.ReadFile("./testdata/simple.png")
	require.NoError(t, err)

	// To Simulate Remote Server
	mux := http.NewServeMux()
	mux.HandleFunc("/image.png", func(w http.ResponseWriter, r *http.Request) {
		bytesReader := bytes.NewReader(dummyImage)
		w.WriteHeader(http.StatusOK)
		bufferedWriter := bufio.NewWriter(w)
		_, err = io.Copy(bufferedWriter, bytesReader)
		require.NoError(t, err)
	})
	client, teardown := httputil.DummyClient(mux)
	return dummyImage, client, teardown
}

func assertImageContent(err error, dummyFs afero.Fs, got *image.Image, t *testing.T, dummyImage []byte) {
	f, err := dummyFs.Open(filepath.Join("./data/images", got.Location))
	assert.NoError(t, err)
	defer f.Close()
	gotContent, err := ioutil.ReadAll(f)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(dummyImage, gotContent))
}

func assertImage(t *testing.T, got *image.Image, size int64, name, ext, userId, desc string) {
	assert.Equal(t, name, got.Name)
	assert.True(t, strings.HasSuffix(got.Location, ext))
	assert.Equal(t, userId, got.UserID)
	assert.NotEmpty(t, got.ImageID)
	assert.NotEmpty(t, got.CreatedAt)
	assert.Equal(t, got.Description, desc)
	assert.Equal(t, size, got.Size)
}
