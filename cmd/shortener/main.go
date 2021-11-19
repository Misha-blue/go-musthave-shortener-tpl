package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var urls = []string{}

var handeledURL = map[string]string{}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getShortenURL(s []string, e string) string {
	for i, a := range s {
		if a == e {
			return fmt.Sprintf("%d", i)
		}
	}
	return ""
}

func HandleURLRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ids := strings.Split(r.URL.Path, "/")
		id := ids[len(ids)-1]
		if handeledURL[id] != "" {
			w.Header().Add("Content-Type", "application/json")
			w.Header().Set("Location", handeledURL[id])
			w.WriteHeader(http.StatusTemporaryRedirect)
			w.Write([]byte(handeledURL[id]))
		} else {
			http.Error(w, "No such id", http.StatusBadRequest)
		}
	case http.MethodPost:
		defer r.Body.Close()
		body, error := io.ReadAll(r.Body)
		if error != nil {
			http.Error(w, error.Error(), 400)
			return
		}
		url := string(body)
		if !contains(urls, url) {
			urls = append(urls, url)
			handeledURL[getShortenURL(urls, url)] = url
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + getShortenURL(urls, url)))
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", HandleURLRequest)

	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
