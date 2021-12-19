package repository

type Storager interface {
	GetAll() (map[string]string, error)
	Add(shortURL string) (string, error)
}
