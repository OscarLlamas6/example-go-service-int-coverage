package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"e2e-test/mocks"
	"e2e-test/src/cache"

	"github.com/julienschmidt/httprouter"
	"github.com/matryer/is"
)

const (
	testKey   = "testKey"
	testValue = "testValue"
)

var errRedis = errors.New("testRedisErr")

func TestGet(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string

		redisGetResp string
		redisGetErr  error

		expectedStatus int
		expectedResp   string
	}{
		{
			name:           "value exists",
			redisGetResp:   testValue,
			expectedStatus: http.StatusOK,
			expectedResp:   testValue,
		},
		{
			name:           "value does not exist",
			redisGetErr:    cache.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedResp:   "key not found\n",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			mockRedis := mocks.NewRedisWrapper(t)

			mockRedis.On("Get", testKey).Return(tc.redisGetResp, tc.redisGetErr)

			srv := &server{Redis: mockRedis}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", testKey), nil)
			w := httptest.NewRecorder()

			router := httprouter.New()
			router.GET(routeKey, srv.HandleGet)
			router.ServeHTTP(w, req)

			is.Equal(tc.expectedStatus, w.Code)
			is.Equal(tc.expectedResp, w.Body.String())
		})
	}
}

func TestSet(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string

		redisSetErr error

		expectedStatus int
	}{
		{
			name:           "success",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "redis error",
			redisSetErr:    errRedis,
			expectedStatus: http.StatusInternalServerError,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			mockRedis := mocks.NewRedisWrapper(t)

			mockRedis.On("Set", testKey, testValue).Return(tc.redisSetErr)

			srv := &server{Redis: mockRedis}

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/%s", testKey, testValue), nil)
			w := httptest.NewRecorder()

			router := httprouter.New()
			router.POST(routeKeyValue, srv.HandleSet)
			router.ServeHTTP(w, req)

			is.Equal(tc.expectedStatus, w.Code)
		})
	}
}

func TestDel(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string

		redisSetErr error

		expectedStatus int
	}{
		{
			name:           "success",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "redis error",
			redisSetErr:    errRedis,
			expectedStatus: http.StatusInternalServerError,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			mockRedis := mocks.NewRedisWrapper(t)

			mockRedis.On("Del", testKey).Return(tc.redisSetErr)

			srv := &server{Redis: mockRedis}

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", testKey), nil)
			w := httptest.NewRecorder()

			router := httprouter.New()
			router.DELETE(routeKey, srv.HandleDel)
			router.ServeHTTP(w, req)

			is.Equal(tc.expectedStatus, w.Code)
		})
	}
}
