package handlers

import (
	"bytes"
	"encoding/json"
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
	r, err := NewRouter()
	require.NoError(t, err)

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

func TestHandleURLJsonPostRequest(t *testing.T) {
	type Request struct {
		URL string `json:"url"`
	}
	type Response struct {
		Result string `json:"result"`
	}

	type want struct {
		contentType string
		statusCode  int
		response    Response
	}
	tests := []struct {
		name    string
		request Request
		want    want
	}{
		{
			name:    "add new Url",
			request: Request{"newUrl"},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				response:    Response{"http://localhost:8080/1"},
			},
		},
		{
			name:    "add existing Url",
			request: Request{"existingUrl"},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				response:    Response{"http://localhost:8080/0"},
			},
		},
	}
	r, err := NewRouter()
	require.NoError(t, err)

	ts := httptest.NewServer(r)
	defer ts.Close()

	existingURL, _ := json.Marshal(Request{"existingUrl"})
	postResp, _ := testRequest(t, ts, http.MethodPost, "/api/shorten", bytes.NewReader(existingURL))
	assert.Equal(t, http.StatusCreated, postResp.StatusCode)
	defer postResp.Body.Close()
	var response Response

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testURL, _ := json.Marshal(tt.request)
			resp, body := testRequest(t, ts, http.MethodPost, "/api/shorten", bytes.NewReader(testURL))
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			_ = json.Unmarshal([]byte(body), &response)
			assert.Equal(t, tt.want.response, response)
		})
	}
}

func TestHandleURLWrongJsonPostRequest(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "add wrong json Url",
			request: "newUrl",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				response:    "invalid character 'e' in literal null (expecting 'u')\n",
			},
		},
	}
	r, err := NewRouter()
	require.NoError(t, err)

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, http.MethodPost, "/api/shorten", bytes.NewReader([]byte(tt.request)))
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, body)
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
	r, err := NewRouter()
	require.NoError(t, err)
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
	r, err := NewRouter()
	require.NoError(t, err)
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

func NewRouter() (chi.Router, error) {
	fmt.Print("test")
	r := chi.NewRouter()

	repository, err := repository.New("fileStorage.txt")
	if err != nil {
		return nil, err
	}

	handler := New(repository, "http://localhost:8080")

	r.Get("/{shortURL}", handler.HandleURLGetRequest)
	r.Post("/", handler.HandleURLPostRequest)
	r.Post("/api/shorten", handler.HandleURLJsonPostRequest)
	return r, nil
}
