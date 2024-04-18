package service

import (
	"app/pkg/logs"
	"app/pkg/smtp"
	"app/storage"
)

type IServiceManager interface {
}

type Service struct {
	store       storage.StoreInterface
	log         logs.LoggerInterface
	filestorage storage.FileStorageInterface
	cache       storage.CacheInterface
	smtp        smtp.SMTPInterface
}

func New(str storage.StoreInterface,
	log logs.LoggerInterface,
	filestorage storage.FileStorageInterface,
	cache storage.CacheInterface,
	stmp smtp.SMTPInterface,
) IServiceManager {
	return &Service{
		store:       str,
		log:         log,
		filestorage: filestorage,
		cache:       cache,
		smtp:        stmp,
	}
}
