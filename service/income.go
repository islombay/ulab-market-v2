package service

import (
	models_v1 "app/api/models/v1"
	"app/pkg/logs"
	"app/storage"
	"context"
)

type incomeService struct {
	store storage.StoreInterface
	log   logs.LoggerInterface
}

func NewIncomeService(store storage.StoreInterface, log logs.LoggerInterface) *incomeService {
	return &incomeService{
		log:   log,
		store: store,
	}
}

// income

func (i *incomeService) Create(ctx context.Context, income models_v1.CreateIncome) (models_v1.CreateIncomeResponse, error) {
	var incomeResponse models_v1.CreateIncomeResponse

	createdIncome, err := i.store.Income().Create(ctx, income)
	if err != nil {
		i.log.Error("error is while creating income", logs.Error(err))
		return models_v1.CreateIncomeResponse{}, err
	}

	for _, incomeProduct := range income.Products {
		createdIncomeProduct, err := i.store.IncomeProduct().CreateIncomeProduct(ctx, models_v1.CreateIncomeProduct{
			IncomeID:     createdIncome.ID,
			ProductID:    incomeProduct.ProductID,
			Quantity:     incomeProduct.Quantity,
			ProductPrice: incomeProduct.ProductPrice,
			TotalPrice:   incomeProduct.TotalPrice,
		})
		if err != nil {
			i.log.Error("error is while creating income product", logs.Error(err))
			return models_v1.CreateIncomeResponse{}, err
		}

		incomeResponse.IncomeProducts = append(incomeResponse.IncomeProducts, createdIncomeProduct)
	}

	incomeResponse.Income = createdIncome

	return incomeResponse, nil
}

func (i *incomeService) GetByID(ctx context.Context, id string) (models_v1.Income, error) {
	return i.store.Income().GetByID(ctx, id)
}

func (i *incomeService) GetList(ctx context.Context, request models_v1.IncomeRequest) (models_v1.IncomeResponse, error) {
	return i.store.Income().GetList(ctx, request)
}

// income_product

func (i *incomeService) GetByIncomeProductID(ctx context.Context, id string) (models_v1.IncomeProduct, error) {
	return i.store.IncomeProduct().GetByIncomeProductID(ctx, id)
}

func (i *incomeService) GetIncomeProductsList(ctx context.Context, request models_v1.IncomeProductRequest) (models_v1.IncomeProductResponse, error) {
	return i.store.IncomeProduct().GetIncomeProductList(ctx, request)
}
