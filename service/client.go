package service

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/helper"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"
	"net/http"
	"time"

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

func (srv *clientService) GetList(ctx context.Context, pagination models.Pagination) (interface{}, *status.Status) {
	model, count, err := srv.store.User().GetClientList(ctx, pagination)
	if err != nil {
		srv.log.Error("could not get list of clients", logs.Error(err))
		return nil, &status.StatusInternal
	}

	var resp_model []models.ClientListAdminPanel

	for _, usr := range model {
		count, err := srv.getOrdersCount(ctx, usr.ID)
		if err != nil {
			return nil, err
		}
		tmp := models.ClientListAdminPanel{}
		if err := helper.Reobject(usr, &tmp, "obj"); err != nil {
			srv.log.Error("could not reobject models.Client to models.ClientListAdminPanel",
				logs.Error(err))
			return nil, &status.StatusInternal
		}

		tmp.OrderCount = count
		resp_model = append(resp_model, tmp)
	}

	return models.Response{
		StatusCode: http.StatusOK,
		Count:      count,
		Data:       resp_model,
	}, nil
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

func (srv *clientService) GetMe(ctx context.Context, user_id string) (interface{}, *status.Status) {
	model, err := srv.store.User().GetClientByID(ctx, user_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &status.StatusUserNotFound
		}

		srv.log.Error("could not find client by id", logs.Error(err), logs.String("uid", user_id))
		return nil, &status.StatusInternal
	}

	if model.DeletedAt != nil {
		return nil, &status.StatusUserNotFound
	}

	return model, nil
}

func (srv *clientService) Update(ctx context.Context, model models.ClientUpdate) (interface{}, *status.Status) {
	if model.Email != nil {
		if *model.Email == "" {
			model.Email = nil
		} else {
			if !helper.IsValidEmail(*model.Email) {
				return nil, &status.StatusBadEmail
			}
		}
	}

	if model.Gender != nil {
		if !(*model.Gender == "male" || *model.Gender == "female") {
			return nil, &status.StatusBadGender
		}
	}

	if model.BirthDate != nil {
		if model.BirthDate.After(time.Now()) {
			return nil, &status.StatusBadDate
		}
	}

	if model.Name != nil {
		if len(*model.Name) < 3 || len(*model.Name) > 20 {
			return nil, &status.StatusNameInvalid
		}
	}

	if model.Surname != nil {
		if len(*model.Surname) < 3 || len(*model.Surname) > 20 {
			return nil, &status.StatusSurnameInvalid
		}
	}

	if err := srv.store.User().UpdateClient(ctx, model); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return nil, &status.StatusAlreadyExists
		}
		srv.log.Error("could not update client", logs.Error(err), logs.String("uid", model.ID))
		return nil, &status.StatusInternal
	}

	return models_v1.Response{Message: "Ok", Code: http.StatusOK}, nil
}
