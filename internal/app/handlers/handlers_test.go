package handlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleURLPostRequest(t *testing.T) {
	type want struct {
		contentType  string
		statusCode   int
		shortenedURL string
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
				shortenedURL: "http://localhost:8080/1",
			},
		},
		{
			name: "add existing Url",
			url:  "existingUrl",
			want: want{
				contentType:  "application/json",
				statusCode:   http.StatusCreated,
				shortenedURL: "http://localhost:8080/0",
			},
		},
	}
	r := NewRouter()

	ts := httptest.NewServer(r)
	defer ts.Close()

	postResp, postBody := testRequest(t, ts, http.MethodPost, "/", bytes.NewReader([]byte("existingUrl")))
	assert.Equal(t, http.StatusCreated, postResp.StatusCode)
	assert.Equal(t, "http://localhost:8080/0", postBody)
	defer postResp.Body.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, http.MethodPost, "/", bytes.NewReader([]byte(tt.url)))
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.shortenedURL, body)
		})
	}
}

func TestHandleURLGetRequest(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		location    string
		body        string
	}
	tests := []struct {
		name  string
		urlID string
		want  want
	}{
		{
			name:  "not existing Url",
			urlID: "1000",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				location:    "",
				body:        "record in storage for your shortUrl wasn't found\n",
			},
		},
		{
			name:  "existing Url",
			urlID: "0",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusTemporaryRedirect,
				location:    "anotherExistingUrl",
				body:        "",
			},
		},
	}
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	postResp, postBody := testRequest(t, ts, http.MethodPost, "/", bytes.NewReader([]byte("anotherExistingUrl")))
	assert.Equal(t, http.StatusCreated, postResp.StatusCode)
	assert.Equal(t, "http://localhost:8080/0", postBody)
	defer postResp.Body.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, http.MethodGet, "/"+tt.urlID, nil)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.location, resp.Header.Get("Location"))
			assert.Equal(t, tt.want.body, body)
		})
	}
}

func TestHandleURLOtherMethodsRequest(t *testing.T) {
	type want struct {
		statusCode int
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
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		want := want{
			statusCode: http.StatusMethodNotAllowed,
		}
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, tt.method, "/", nil)
			defer resp.Body.Close()
			assert.Equal(t, want.statusCode, resp.StatusCode)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func NewRouter() chi.Router {
	fmt.Print("test")
	r := chi.NewRouter()

	handler := New(repository.New())

	r.Get("/{shortURL}", handler.HandleURLGetRequest)
	r.Post("/", handler.HandleURLPostRequest)
	return r
}
