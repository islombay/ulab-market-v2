package service

import (
	"app/pkg/logs"
	"app/pkg/smtp"
	"app/storage"
)

type IServiceManager interface {
	Order() OrderService
	Favourite() *FavouriteService
}

type Service struct {
	order     OrderService
	favourite *FavouriteService
}

func New(str storage.StoreInterface, log logs.LoggerInterface, filestorage storage.FileStorageInterface, cache storage.CacheInterface, stmp smtp.SMTPInterface) IServiceManager {
	srv := Service{}

	srv.order = NewOrderService(str, log, filestorage)
	srv.favourite = NewFavouriteService(str, log)

	return &srv
}

func (s *Service) Order() OrderService {
	return s.order
}

func (s *Service) Favourite() *FavouriteService {
	return s.favourite
}

var (
	ContextUserID = "user_id"
)
