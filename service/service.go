package service

import (
	"app/pkg/logs"
	"app/pkg/smtp"
	"app/storage"
)

type IServiceManager interface {
	Order() OrderService
	Store() *storeService
	Favourite() *FavouriteService
	Income() *IncomeService
}

type Service struct {
	order        OrderService
	favourite    *FavouriteService
	storeService *storeService

	income *IncomeService
}

func New(str storage.StoreInterface, log logs.LoggerInterface, filestorage storage.FileStorageInterface, cache storage.CacheInterface, stmp smtp.SMTPInterface) IServiceManager {
	srv := Service{}

	srv.order = NewOrderService(str, log, filestorage)
	srv.storeService = NewStoreService(str, log)
	srv.favourite = NewFavouriteService(str, log)
	srv.income = NewIncomeService(str, log)

	return &srv
}

func (s *Service) Order() OrderService {
	return s.order
}

func (s *Service) Store() *storeService {
	return s.storeService
}

func (s *Service) Favourite() *FavouriteService {
	return s.favourite
}

func (s *Service) Income() *IncomeService {
	return s.income
}

var (
	ContextUserID = "user_id"
)
