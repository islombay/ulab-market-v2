package storage

import (
	"mime/multipart"
)

type Folder string

type FileStorageInterface interface {
	Create(model *multipart.FileHeader, imageFolder Folder, id string) (string, error)
	GetURL(id string) string
	DeleteFile(path string) error
}
