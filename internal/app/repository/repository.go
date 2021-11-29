package repository

import (
	"errors"
	"fmt"
)

type Repositorier map[string]string

type Repositorer interface {
	Store(url string) (string, error)
	Load(shortURL string) (string, error)
}

func New() Repositorier {
	return Repositorier{}
}

func (repo Repositorier) Store(url string) (string, error) {
	err := error(nil)

	shortURL := findShortURL(repo, url)

	if shortURL == "" {
		shortURL = generateShortURL(repo)
		repo[shortURL] = url
	}

	return shortURL, err
}

func (repo Repositorier) Load(shortURL string) (string, error) {
	err := error(nil)

	url := repo[shortURL]

	if url == "" {
		err = errors.New("record in storage for your shortUrl wasn't found")
	}

	return url, err
}

func findShortURL(repo Repositorier, url string) string {
	for key, value := range repo {
		if value == url {
			return key
		}
	}
	return ""
}

func generateShortURL(repo Repositorier) string {
	return fmt.Sprintf("%d", len(repo))
}
