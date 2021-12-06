package repository

type Storage interface {
	GetAll() (map[string]string, error)
	Add(shortURL string) (string, error)
}
