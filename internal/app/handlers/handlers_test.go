package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleURLRequest_post(t *testing.T) {
	type want struct {
		contentType  string
		statusCode   int
		shortenedUrl string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "add new Url",
			url:  "newUrl",
			want: want{
				contentType:  "application/json",
				statusCode:   http.StatusCreated,
				shortenedUrl: "http://localhost:8080/1",
			},
		},
		{
			name: "add existing Url",
			url:  "existingUrl",
			want: want{
				contentType:  "application/json",
				statusCode:   http.StatusOK,
				shortenedUrl: "http://localhost:8080/0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1 := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("existingUrl")))
			w1 := httptest.NewRecorder()
			h1 := http.HandlerFunc(HandleURLRequest)
			h1.ServeHTTP(w1, r1)

			body := bytes.NewReader([]byte(tt.url))
			request := httptest.NewRequest(http.MethodPost, "/", body)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandleURLRequest)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			resultBody, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.shortenedUrl, string(resultBody))
		})
	}
}

func TestHandleURLRequest_get(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		location    string
		body        string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "existing Url",
			url:  "0",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusTemporaryRedirect,
				location:    "existingUrl",
				body:        "",
			},
		},
		{
			name: "not existing Url",
			url:  "1000",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				location:    "",
				body:        "Invalid shortened url id.\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1 := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("existingUrl")))
			w1 := httptest.NewRecorder()
			h1 := http.HandlerFunc(HandleURLRequest)
			h1.ServeHTTP(w1, r1)

			request := httptest.NewRequest(http.MethodGet, "/"+tt.url, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandleURLRequest)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			resultBody, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, string(resultBody))
		})
	}
}

func TestHandleURLRequest_othermethods(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	tests := []struct {
		name   string
		method string
	}{
		{
			name:   "head method",
			method: http.MethodHead,
		},
		{
			name:   "patch method",
			method: http.MethodPatch,
		},
		{
			name:   "delete method",
			method: http.MethodDelete,
		},
		{
			name:   "put method",
			method: http.MethodPut,
		},
		{
			name:   "trace method",
			method: http.MethodTrace,
		},
		{
			name:   "options method",
			method: http.MethodOptions,
		},
		{
			name:   "connect method",
			method: http.MethodConnect,
		},
	}
	for _, tt := range tests {
		want := want{
			contentType: "text/plain; charset=utf-8",
			statusCode:  http.StatusBadRequest,
			body:        "Unsupported method",
		}
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandleURLRequest)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, want.statusCode, result.StatusCode)
			assert.Equal(t, want.contentType, result.Header.Get("Content-Type"))
			resultBody, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Contains(t, string(resultBody), want.body)
		})
	}
}
