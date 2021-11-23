package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
)

var handeledURLs = map[string]string{}

func contains(storage map[string]string, element string) bool {
	for _, value := range storage {
		if value == element {
			return true
		}
	}
	return false
}

func findShornedURL(storage map[string]string, element string) string {
	for key, value := range storage {
		if value == element {
			return key
		}
	}
	return ""
}

func getShortenURL(storage map[string]string, e string) string {
	return fmt.Sprintf("%d", len(storage))
}

func HandleURLPostRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, error := io.ReadAll(r.Body)
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}
	url := string(body)
	w.Header().Add("Content-Type", "application/json")
	if !contains(handeledURLs, url) {
		shortID := getShortenURL(handeledURLs, url)
		handeledURLs[shortID] = url
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + shortID))
	} else {
		shortID := findShornedURL(handeledURLs, url)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("http://localhost:8080/" + shortID))
	}
}

func HandleURLGetRequest(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "urlID")
	if handeledURLs[urlID] != "" {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Location", handeledURLs[urlID])
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Invalid shortened url id.", http.StatusBadRequest)
	}
}
