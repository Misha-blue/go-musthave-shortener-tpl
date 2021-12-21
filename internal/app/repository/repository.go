package repository

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Repository struct {
	filePath string
}

type Repositorer interface {
	Store(url string) (string, error)
	Load(shortURL string) (string, error)
}

func New(filePath string) (*Repository, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return &Repository{filePath: filePath}, nil
}

func (repo Repository) Store(url string) (string, error) {
	err := error(nil)

	urls, e := getAll(repo.filePath)
	if e != nil {
		return "", e
	}

	shortURL := findShortURL(urls, url)

	if shortURL == "" {
		shortURL = generateShortURL(urls)
		_, err = add(repo.filePath, shortURL, url)
	}

	return shortURL, err
}

func (repo Repository) Load(shortURL string) (string, error) {
	err := error(nil)

	urls, e := getAll(repo.filePath)
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

func getAll(filePath string) (map[string]string, error) {
	storage := make(map[string]string)
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		record := strings.Split(scanner.Text(), ";")
		storage[record[0]] = record[1]
	}

	return storage, nil
}

func add(filePath string, shortURL string, originURL string) (int, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Write([]byte(shortURL + ";" + originURL + "\n"))
}
