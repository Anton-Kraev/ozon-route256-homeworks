package storage

type Storage struct {
	fileName string
}

func NewStorage(fileName string) Storage {
	return Storage{fileName: fileName}
}
