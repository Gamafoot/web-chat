package storage

type Storage struct {
	User UserStorage
}

func NewStorage(user UserStorage) Storage {
	return Storage{
		User: user,
	}
}
