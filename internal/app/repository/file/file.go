package file

import (
	"bufio"
	"os"
	"strings"
)

type FileStorage struct {
	file    *os.File
	scanner *bufio.Scanner
}

func New(filePath string) (*FileStorage, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &FileStorage{
		file:    file,
		scanner: bufio.NewScanner(file)}, nil
}

func (s *FileStorage) GetAll() (map[string]string, error) {
	result := make(map[string]string)

	if !s.scanner.Scan() {
		return nil, s.scanner.Err()
	}

	for s.scanner.Scan() {
		record := strings.Split(s.scanner.Text(), ";")
		result[record[0]] = record[1]
	}

	return result, nil
}

func (s *FileStorage) Add(shortURL string, originURL string) (int, error) {
	return s.file.Write([]byte(shortURL + ";" + originURL))
}

func (s *FileStorage) Close() error {
	return s.file.Close()
}
