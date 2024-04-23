package service

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type FavouriteService struct {
	store storage.StoreInterface
	log   logs.LoggerInterface
}

func NewFavouriteService(
	store storage.StoreInterface,
	log logs.LoggerInterface,
) *FavouriteService {
	return &FavouriteService{
		store: store,
		log:   log,
	}
}

func (srv *FavouriteService) AddFavourite(ctx context.Context, m models_v1.AddToFavourite) (interface{}, *status.Status) {
	userID, ok := ctx.Value(ContextUserID).(string)
	if !ok {
		srv.log.Error("invalid user id")
		return nil, &status.StatusUnauthorized
	}

	if _, err := srv.store.Product().GetByID(ctx, m.ProductID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &status.StatusProductNotFount
		}
		srv.log.Error("could not get product by id", logs.Error(err),
			logs.String("pid", m.ProductID),
		)
		return nil, &status.StatusInternal
	}

	if _, err := srv.store.Favourite().Get(ctx, userID, m.ProductID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			srv.log.Error("could not get favourite",
				logs.Error(err),
				logs.String("uid", userID),
				logs.String("pid", m.ProductID),
			)
			return nil, &status.StatusInternal
		}
	} else {
		return nil, &status.StatusAlreadyExists
	}

	if err := srv.store.Favourite().Create(ctx, userID, m.ProductID); err != nil {
		srv.log.Error("could not create favourite in db",
			logs.Error(err),
			logs.String("uid", userID),
			logs.String("pid", m.ProductID),
		)
		return nil, &status.StatusInternal
	}

	return models.FavouriteModel{
		ProductID: m.ProductID,
	}, nil
}
