package service

import (
	"app/api/status"
	"app/pkg/logs"
	"app/storage"
	"context"
)

type clientService struct {
	store storage.StoreInterface
	log   logs.LoggerInterface
}

func NewClientService(
	store storage.StoreInterface,
	log logs.LoggerInterface,
) *clientService {
	return &clientService{
		store: store,
		log:   log,
	}
}

func (srv *clientService) GetList(ctx context.Context) (interface{}, *status.Status) {
	model, err := srv.store.User().GetClientList(ctx)
	if err != nil {
		srv.log.Error("could not get list of clients", logs.Error(err))
		return nil, &status.StatusInternal
	}

	return model, nil
}
