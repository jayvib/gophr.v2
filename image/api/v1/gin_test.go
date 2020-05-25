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
		UserID: userutil.GenerateID(),
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
		"name": img.Name,
		"description": img.Description,
		"username": "gopher",
	})

	resp := httputil.PerformRequest(e, http.MethodPost, "/image/file", body, func(r *http.Request){
		r.Header.Add("Content-Type", contentType)
	})

	assert.Equal(t, http.StatusCreated, resp.Code)
	assertImageFromResponse(t, resp)

	userService.AssertExpectations(t)
	repo.AssertExpectations(t)
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
	for k, v := range fields {
		err = multipartWriter.WriteField(k, v)
		require.NoError(t, err)
	}
	return body, multipartWriter.FormDataContentType()
}
