package service

import (
	models_v1 "app/api/models/v1"
	"app/pkg/logs"
	"app/storage"
	"context"
)

type storeService struct {
	store storage.StoreInterface
	log   logs.LoggerInterface
}

func NewStoreService(store storage.StoreInterface, log logs.LoggerInterface) storeService {
	return storeService{
		log:   log,
		store: store,
	}
}

func (s *storeService) Create(ctx context.Context, createStorage models_v1.CreateStorage) stri {

}
