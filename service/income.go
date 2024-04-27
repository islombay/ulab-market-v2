package service

import (
	models_v1 "app/api/models/v1"
	"app/pkg/logs"
	"app/storage"
	"context"
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

func (i *IncomeService) Create(ctx context.Context, income models_v1.CreateIncome) (models_v1.CreateIncomeResponse, error) {
	var incomeResponse models_v1.CreateIncomeResponse

	createdIncome, err := i.store.Income().Create(ctx, income)
	if err != nil {
		i.log.Error("error is while creating income", logs.Error(err))
		return models_v1.CreateIncomeResponse{}, err
	}

	// o'chirildi sababi postgres'ni o'zida transaction orqali income_products qo'shilyabti
	// for _, incomeProduct := range income.Products {
	// 	createdIncomeProduct, err := i.store.Income().CreateIncomeProduct(ctx, models_v1.CreateIncomeProduct{
	// 		IncomeID:     createdIncome.ID,
	// 		ProductID:    incomeProduct.ProductID,
	// 		Quantity:     incomeProduct.Quantity,
	// 		ProductPrice: incomeProduct.ProductPrice,
	// 		TotalPrice:   incomeProduct.TotalPrice,
	// 	})
	// 	if err != nil {
	// 		i.log.Error("error is while creating income product", logs.Error(err))
	// 		return models_v1.CreateIncomeResponse{}, err
	// 	}

	// 	incomeResponse.IncomeProducts = append(incomeResponse.IncomeProducts, createdIncomeProduct)
	// }

	incomeResponse.Income = createdIncome

	return incomeResponse, nil
}

func (i *IncomeService) GetByID(ctx context.Context, id string) (models_v1.Income, error) {
	return i.store.Income().GetByID(ctx, id)
}

func (i *IncomeService) GetList(ctx context.Context, request models_v1.IncomeRequest) (models_v1.IncomeResponse, error) {
	return i.store.Income().GetList(ctx, request)
}

// income_product

func (i *IncomeService) GetByIncomeProductID(ctx context.Context, id string) (models_v1.IncomeProduct, error) {
	return i.store.Income().GetByIncomeProductID(ctx, id)
}

func (i *IncomeService) GetIncomeProductsList(ctx context.Context, request models_v1.IncomeProductRequest) (models_v1.IncomeProductResponse, error) {
	return i.store.Income().GetIncomeProductsList(ctx, request)
}
