package repository

import (
	"errors"
	"fmt"
)

type Storage map[string]string

var storage = Storage{}

type Repositorer interface {
	Store(url string) (string, error)
	Load(shortUrl string) (string, error)
}

func Store(url string) (string, error) {
	err := error(nil)

	shortURL := findShornedURL(storage, url)

	if shortURL == "" {
		shortURL = generateShortenURL(storage)
		storage[shortURL] = url
	}

	return shortURL, err
}

func Load(shortUrl string) (string, error) {
	err := error(nil)

	url := storage[shortUrl]

	if url == "" {
		err = errors.New("record in storage for your shortUrl wasn't found")
	}

	return storage[shortUrl], err
}

func findShornedURL(storage Storage, element string) string {
	for key, value := range storage {
		if value == element {
			return key
		}
	}
	return ""
}

func generateShortenURL(storage Storage) string {
	return fmt.Sprintf("%d", len(storage))
}
