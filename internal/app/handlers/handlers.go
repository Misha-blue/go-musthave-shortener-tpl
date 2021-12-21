package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository"
	"github.com/go-chi/chi"
)

type Handler struct {
	repository *repository.Repository
	baseURL    string
}

type URLPostRequest struct {
	URL string `json:"url"`
}

type URLResponseRequest struct {
	Result string `json:"result"`
}

func New(repository *repository.Repository, baseURL string) *Handler {
	return &Handler{repository: repository, baseURL: baseURL}
}

func (handler *Handler) HandleURLPostRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, e := io.ReadAll(r.Body)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := handler.repository.Store(string(body))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(handler.baseURL + "/" + shortURL))
}

func (handler *Handler) HandleURLJsonPostRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request := URLPostRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := handler.repository.Store(request.URL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := URLResponseRequest{handler.baseURL + "/" + shortURL}

	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.Encode(response)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

func (handler *Handler) HandleURLGetRequest(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	url, err := handler.repository.Load(shortURL)

	if err == nil {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
