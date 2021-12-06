package repository

import (
	"errors"
	"fmt"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository/file"
)

type Repository struct {
	storage *file.FileStorage
}

type Repositorer interface {
	Store(url string) (string, error)
	Load(shortURL string) (string, error)
}

func New(storage *file.FileStorage) Repository {
	return Repository{storage: storage}
}

func (repo Repository) Store(url string) (string, error) {
	err := error(nil)

	urls, e := repo.storage.GetAll()
	if e != nil {
		return "", e
	}

	shortURL := findShortURL(urls, url)

	if shortURL == "" {
		shortURL = generateShortURL(urls)
		_, err = repo.storage.Add(shortURL, url)
	}

	return shortURL, err
}

func (repo Repository) Load(shortURL string) (string, error) {
	err := error(nil)

	urls, e := repo.storage.GetAll()
	if e != nil {
		return "", e
	}

	url := urls[shortURL]

	if url == "" {
		err = errors.New("record in storage for your shortUrl wasn't found")
	}

	return url, err
}

func findShortURL(s map[string]string, url string) string {
	for key, value := range s {
		if value == url {
			return key
		}
	}
	return ""
}

func generateShortURL(s map[string]string) string {
	return fmt.Sprintf("%d", len(s))
}
