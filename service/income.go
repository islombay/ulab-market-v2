package service

import (
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type IncomeService struct {
	store storage.StoreInterface
	log   logs.LoggerInterface
}

func NewIncomeService(store storage.StoreInterface, log logs.LoggerInterface) *IncomeService {
	return &IncomeService{
		log:   log,
		store: store,
	}
}

// income

func (i *IncomeService) Create(ctx context.Context, income models_v1.CreateIncome) (interface{}, *status.Status) {

	createdIncome, err := i.store.Income().Create(ctx, income)
	if err != nil {
		i.log.Error("error is while creating income", logs.Error(err))
		return models_v1.CreateIncomeResponse{}, &status.StatusInternal
	}
	return createdIncome, nil
}

func (i *IncomeService) GetByID(ctx context.Context, id string) (*models_v1.Income, *status.Status) {
	income, err := i.store.Income().GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &status.StatusNotFound
		}
		i.log.Error("could not find the income by id", logs.Error(err),
			logs.String("income id", id))
		return nil, &status.StatusInternal
	}

	income_products, err := i.store.Income().GetProductsByIncomeID(ctx, income.ID)
	if err != nil {
		i.log.Error("could not get products of income",
			logs.Error(err),
			logs.String("income_id", id))
	}

	income.Products = income_products

	return &income, nil
}

func (i *IncomeService) GetList(ctx context.Context, request models_v1.IncomeRequest) (interface{}, *status.Status) {
	if request.Limit == 0 {
		request.Limit = 10
	}
	if request.Page == 0 {
		request.Page = 1
	}
	res, err := i.store.Income().GetList(ctx, request)
	if err != nil {
		i.log.Error("could not get all incomes", logs.Error(err))
		return nil, &status.StatusInternal
	}

	return res, nil
}

// income_product

// func (i *IncomeService) GetByIncomeProductID(ctx context.Context, id string) (models_v1.IncomeProduct, error) {
// 	return i.store.Income().GetByIncomeProductID(ctx, id)
// }

func (i *IncomeService) GetIncomeProductsList(ctx context.Context, request models_v1.IncomeProductRequest) (models_v1.IncomeProductResponse, error) {
	return i.store.Income().GetIncomeProductsList(ctx, request)
}
