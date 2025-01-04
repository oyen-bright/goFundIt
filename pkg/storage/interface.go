package storage

type Storage interface {
	UploadFile(file, folderPath string) (url, id string, err error)
	DeleteFile(id string) error
}
