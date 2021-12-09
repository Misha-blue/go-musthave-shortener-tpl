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
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return &FileStorage{
		filePath: filePath}, nil
}

func (s *FileStorage) GetAll() (map[string]string, error) {
	storage := make(map[string]string)
	file, err := os.Open(s.filePath)

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

func (s *FileStorage) Add(shortURL string, originURL string) (int, error) {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Write([]byte(shortURL + ";" + originURL + "\n"))
}
