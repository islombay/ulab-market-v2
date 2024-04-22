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

func NewStoreService(store storage.StoreInterface, log logs.LoggerInterface) *storeService {
	return &storeService{
		log:   log,
		store: store,
	}
}

func (s *storeService) Create(ctx context.Context, createStorage models_v1.CreateStorage) (models_v1.Storage, error) {
	return s.store.Storage().Create(ctx, createStorage)
}

func (s *storeService) GetByID(ctx context.Context, id string) (models_v1.Storage, error) {
	return s.store.Storage().GetByID(ctx, id)
}
func (s *storeService) GetList(ctx context.Context, request models_v1.StorageRequest) (models_v1.StorageResponse, error) {
	return s.store.Storage().GetList(ctx, request)
}
func (s *storeService) Update(ctx context.Context, updateStorage models_v1.UpdateStorage) (models_v1.Storage, error) {
	return s.store.Storage().Update(ctx, updateStorage)
}
func (s *storeService) Delete(ctx context.Context, id string) error {
	return s.store.Storage().Delete(ctx, id)
}
