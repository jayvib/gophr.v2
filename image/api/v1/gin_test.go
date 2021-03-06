package v1

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophr.v2/http/httputil"
	"gophr.v2/user"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"gophr.v2/image"
	"gophr.v2/image/imageutil"
	"gophr.v2/image/mocks"
	"gophr.v2/image/service"
	usermocks "gophr.v2/user/mocks"
	"gophr.v2/user/userutil"

	"github.com/gin-gonic/gin"
	"github.com/jayvib/golog"
	gophrtesting "gophr.v2/testing"
)

var debug = flag.Bool("debug", false, "Debugging")

func TestMain(m *testing.M) {
	flag.Parse()
	if *debug {
		golog.Info("Debugging Mode!")
		golog.SetLevel(golog.DebugLevel)
	}
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestCreateImageFromFile(t *testing.T) {
	img := &image.Image{
		UserID:      userutil.GenerateID(),
		ImageID:     imageutil.GenerateID(),
		Name:        "Unit Test Image",
		Size:        12345,
		Description: "This is a unit test",
	}

	usr := &user.User{
		UserID:   userutil.GenerateID(),
		Username: "luffy.monkey",
	}

	repo := new(mocks.Repository)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*image.Image")).Return(nil).Once()

	userService := new(usermocks.Service)
	userService.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(usr, nil).Once()

	svc := service.New(repo, afero.NewMemMapFs(), nil)

	e := gin.Default()
	RegisterRoutes(e, svc, userService)

	// Create a multipart
	body, contentType := createMultipartBody(t, "testdata/simple.png", map[string]string{
		"name":        img.Name,
		"description": img.Description,
		"username":    "gopher",
	})

	resp := httputil.PerformRequest(e, http.MethodPost, "/image/file", body, func(r *http.Request) {
		r.Header.Add("Content-Type", contentType)
	})

	assert.Equal(t, http.StatusCreated, resp.Code)
	assertImageFromResponse(t, resp)

	userService.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCreateImageFromURL(t *testing.T) {
	img := &image.Image{
		UserID:      userutil.GenerateID(),
		ImageID:     imageutil.GenerateID(),
		Name:        "Unit Test Image",
		Size:        12345,
		Description: "This is a unit test",
	}

	usr := &user.User{
		UserID:   userutil.GenerateID(),
		Username: "luffy.monkey",
	}

	repo := new(mocks.Repository)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*image.Image")).Return(nil).Once()

	userService := new(usermocks.Service)
	userService.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(usr, nil).Once()

	// Create Stub Remote Server
	testFile, err := os.Open("testdata/simple.png")
	require.NoError(t, err)
	defer testFile.Close()

	isStubRemoteHandlerCalled := false
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(w, testFile)
		require.NoError(t, err)
		isStubRemoteHandlerCalled = true
	})

	stubClient, teardown := gophrtesting.RemoteServerStub(h)
	defer teardown()
	svc := service.New(repo, afero.NewMemMapFs(), stubClient)

	e := gin.Default()
	RegisterRoutes(e, svc, userService)

	// Create a multipart form
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	addFieldsToMultipartWriter(t, mw, map[string]string{
		"name":        img.Name,
		"description": img.Description,
		"username":    "gopher",
		"url":         "http://testing.net/simple.png",
	})
	err = mw.Close()
	require.NoError(t, err)

	// Do a request
	resp := httputil.PerformRequest(e, http.MethodPost, "/image/url", body, func(r *http.Request) {
		r.Header.Add("Content-Type", mw.FormDataContentType())
	})

	assert.Equal(t, http.StatusCreated, resp.Code)
	assertImageFromResponse(t, resp)

	userService.AssertExpectations(t)
	repo.AssertExpectations(t)
	assert.True(t, isStubRemoteHandlerCalled)
}

func TestFind(t *testing.T) {
	want := &image.Image{
		UserID:      userutil.GenerateID(),
		ImageID:     imageutil.GenerateID(),
		Name:        "Unit Test Image",
		Size:        12345,
		Description: "This is a unit test",
	}

	svc := new(mocks.Service)
	svc.On("Find", mock.Anything, mock.AnythingOfType("string")).Return(want, nil).Once()

	e := gin.Default()
	RegisterRoutes(e, svc, nil)

	resp := httputil.PerformRequest(e, http.MethodGet, "/image/id/"+want.UserID, nil)

	require.Equal(t, http.StatusFound, resp.Code)

	var got image.Image
	err := json.NewDecoder(resp.Body).Decode(&got)
	require.NoError(t, err)

	assert.Equal(t, want, &got)
}

func TestFindAllByUser(t *testing.T) {
	userid := "1234abcde"
	images := []*image.Image{
		{
			ID:      1,
			UserID:  userid,
			ImageID: imageutil.GenerateID(),
			Name:    "bacon",
		},
		{
			ID:      2,
			UserID:  userid,
			ImageID: imageutil.GenerateID(),
			Name:    "cheeze",
		},
	}

	svc := new(mocks.Service)
	svc.On("FindAllByUser", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int")).Return(images, nil)

	e := gin.Default()
	RegisterRoutes(e, svc, nil)

	resp := httputil.PerformRequest(e, http.MethodGet, "/image/userid/1234abc?offset=0", nil)

	assert.Equal(t, http.StatusOK, resp.Code)

	got := make([]*user.User, 0)
	err := json.NewDecoder(resp.Body).Decode(&got)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestFindAll(t *testing.T) {
	images := []*image.Image{
		{
			ID:      1,
			UserID:  userutil.GenerateID(),
			ImageID: imageutil.GenerateID(),
			Name:    "bacon",
		},
		{
			ID:      2,
			UserID:  userutil.GenerateID(),
			ImageID: imageutil.GenerateID(),
			Name:    "cheeze",
		},
	}

	svc := new(mocks.Service)
	svc.On("FindAll", mock.Anything, mock.AnythingOfType("int")).Return(images, nil).Once()

	e := gin.Default()
	RegisterRoutes(e, svc, nil)
	resp := httputil.PerformRequest(e, http.MethodGet, "/image?offset=0", nil)

	assert.Equal(t, http.StatusOK, resp.Code)

	got := make([]*user.User, 0)
	err := json.NewDecoder(resp.Body).Decode(&got)
	require.NoError(t, err)
	svc.AssertExpectations(t)
}

func assertImageFromResponse(t *testing.T, resp *httptest.ResponseRecorder) {
	var gotImg image.Image
	err := json.NewDecoder(resp.Body).Decode(&gotImg)
	require.NoError(t, err)
	assert.NotEmpty(t, gotImg.ImageID)
	assert.NotEmpty(t, gotImg.CreatedAt)
	assert.NotEmpty(t, gotImg.UserID)
}

func createMultipartBody(t *testing.T, filename string, fields map[string]string) (reader io.Reader, contentType string) {
	t.Helper()
	file, err := os.Open(filename)
	require.NoError(t, err)
	defer file.Close()
	fileContent, err := ioutil.ReadAll(file)
	require.NoError(t, err)
	fi, err := os.Stat(filename)
	require.NoError(t, err)
	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)
	part, err := multipartWriter.CreateFormFile("file", fi.Name())
	require.NoError(t, err)
	_, err = part.Write(fileContent)
	require.NoError(t, err)
	defer multipartWriter.Close()
	// Add some metadata
	addFieldsToMultipartWriter(t, multipartWriter, fields)
	return body, multipartWriter.FormDataContentType()
}

func addFieldsToMultipartWriter(t *testing.T, w *multipart.Writer, fields map[string]string) {
	for k, v := range fields {
		err := w.WriteField(k, v)
		require.NoError(t, err)
	}
}
