package service

import (
	"app/api/models"
	"app/api/status"
	"app/pkg/helper"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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

	var resp_model []models.ClientSwagger

	for _, usr := range model {
		count, err := srv.getOrdersCount(ctx, usr.ID)
		if err != nil {
			return nil, err
		}
		tmp := models.ClientSwagger{}
		if err := helper.Reobject(usr, &tmp, "obj"); err != nil {
			srv.log.Error("could not reobject models.Client to models.ClientSwagger",
				logs.Error(err))
			return nil, &status.StatusInternal
		}

		tmp.OrderCount = count
		resp_model = append(resp_model, tmp)
	}

	return resp_model, nil
}

func (srv *clientService) getOrdersCount(ctx context.Context, user_id string) (int, *status.Status) {
	count, err := srv.store.Order().OrdersCount(ctx, user_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, &status.StatusNotFound
		}
		srv.log.Error("could not get order_count for client", logs.Error(err),
			logs.String("uid", user_id))
		return 0, &status.StatusInternal
	}

	return count, nil
}
