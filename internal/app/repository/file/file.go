package file

import (
	"bufio"
	"os"
	"strings"
)

type FileStorage struct {
	filePath string
}

func New(filePath string) (*FileStorage, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0644)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	return &FileStorage{
		filePath: filePath}, nil
}

func (s *FileStorage) GetAll() (map[string]string, error) {
	storage := make(map[string]string)
	file, err := os.Open(s.filePath)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		record := strings.Split(scanner.Text(), ";")
		storage[record[0]] = record[1]
	}

	return storage, nil
}

func (s *FileStorage) Add(shortURL string, originURL string) (int, error) {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Write([]byte(shortURL + ";" + originURL + "\n"))
}
