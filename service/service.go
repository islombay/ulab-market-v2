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
	Courier() *courierService
	Client() *clientService
	Notify() *notificationService
}

type Service struct {
	order        OrderService
	favourite    *FavouriteService
	storeService *storeService

	income       *IncomeService
	courier      *courierService
	client       *clientService
	notification *notificationService
}

func New(str storage.StoreInterface, log logs.LoggerInterface, filestorage storage.FileStorageInterface, cache storage.CacheInterface, stmp smtp.SMTPInterface) IServiceManager {
	srv := Service{}

	srv.notification = NewNotificationService(log)
	srv.order = NewOrderService(str, log, filestorage, srv.notification)
	srv.storeService = NewStoreService(str, log)
	srv.favourite = NewFavouriteService(str, log)
	srv.income = NewIncomeService(str, log)
	srv.courier = NewCourierService(str, log)
	srv.client = NewClientService(str, log)

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

func (s *Service) Client() *clientService {
	return s.client
}

func (s *Service) Income() *IncomeService {
	return s.income
}

func (s *Service) Courier() *courierService {
	return s.courier
}

func (s *Service) Notify() *notificationService {
	return s.notification
}

var (
	ContextUserID = "user_id"
)
