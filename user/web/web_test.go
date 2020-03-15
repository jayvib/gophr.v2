package web

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/jayvib/golog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gophr.v2/user"
	"gophr.v2/user/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var debugTest = flag.Bool("debug", false, "debugging")

type responseTest struct {
	Message string     `json:"message"`
	Data    *user.User `json:"data"`
	Method  string     `json:"method"`
}

func TestMain(m *testing.M) {
	flag.Parse()
	if *debugTest {
		golog.SetLevel(golog.DebugLevel)
	}
	os.Exit(m.Run())
}

func TestHandle_GetByID(t *testing.T) {

	t.Run("OK", func(t *testing.T) {

		service := new(mocks.Service)

		mockReturn := &user.User{
			ID:       1,
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
			Password: "1234567890",
		}

		want := responseTest{
			Message: "OK",
			Data:    mockReturn,
			Method:  http.MethodGet,
		}

		service.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(mockReturn, nil).Once()

		param := &Parameters{
			UserService: service,
		}

		h := New(param)

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", h.GetByID).Methods(http.MethodGet)

		w := performRequest(router, http.MethodGet, "/users/1", nil)
		assertResponse(t, want, w)
		service.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		service := new(mocks.Service)

		mockReturn := &user.User{
			ID:       1,
			Username: "luffy.monkey",
			Email:    "luffy.monkey@gmail.com",
			Password: "1234567890",
		}

		want := responseTest{
			Message: "OK",
			Data:    mockReturn,
			Method:  http.MethodGet,
		}

		service.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(mockReturn, nil).Once()

		param := &Parameters{
			UserService: service,
		}

		h := New(param)

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", h.GetByID).Methods(http.MethodGet)

		w := performRequest(router, http.MethodGet, "/users/mock1", nil)
		assertResponse(t, want, w)
		service.AssertExpectations(t)

	})
}

func TestHandler_GetByEmail(t *testing.T) {
	svc := new(mocks.Service)
	mockReturn := &user.User{
		ID:       1,
		Username: "luffy.monkey",
		Email:    "luffy.monkey@gmail.com",
		Password: "1234567890",
	}
	_ = mockReturn

	svc.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(mockReturn, nil).Once()
	param := &Parameters{
		UserService: svc,
	}
	h := New(param)
	router := mux.NewRouter()
	router.HandleFunc("/user/email/{email}", h.GetByEmail).Methods(http.MethodGet)
	resp := performRequest(router, http.MethodGet, "/user/email/luffy.monkey%40gmail.com", nil)
	assert.Equal(t, http.StatusOK, resp.Code)
	want := responseTest{
		Data:    mockReturn,
		Message: "OK",
		Method:  http.MethodGet,
	}
	assertResponse(t, want, resp)
}

func assertResponse(t *testing.T, want responseTest, w *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, w.Code)
	var got responseTest
	err := json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func assertGetByIDNotFound(t *testing.T, want responseTest, w *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusNotFound, w.Code)

	var got responseTest

	err := json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)

	assert.Equal(t, want, got)

}

func performRequest(h http.Handler, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}
