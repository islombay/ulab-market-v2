package service

import (
	"app/pkg/logs"
	"app/pkg/smtp"
	"app/storage"
)

type IServiceManager interface {
	Order() OrderService
}

type Service struct {
	order OrderService
}

func New(str storage.StoreInterface, log logs.LoggerInterface, filestorage storage.FileStorageInterface, cache storage.CacheInterface, stmp smtp.SMTPInterface) IServiceManager {
	srv := Service{}

	srv.order = NewOrderService(str, log, filestorage)

	return &srv
}

func (s *Service) Order() OrderService {
	return s.order
}

var (
	ContextUserID = "user_id"
)
