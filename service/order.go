package service

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type OrderService struct {
	store       storage.StoreInterface
	log         logs.LoggerInterface
	filestorage storage.FileStorageInterface
}

func NewOrderService(store storage.StoreInterface,
	log logs.LoggerInterface,
	filestorage storage.FileStorageInterface,
) OrderService {
	return OrderService{
		log:         log,
		store:       store,
		filestorage: filestorage,
	}
}

const (
	OrderStatusInProcess = "in_process"
	OrderStatusFinished  = "finished"
	OrderStatusCanceled  = "canceled"
)

func (srv OrderService) CreateOrder(ctx context.Context, order models_v1.CreateOrder) (interface{}, *status.Status) {
	userID, ok := ctx.Value(ContextUserID).(string)
	if !ok {
		srv.log.Error("invalid user id")
		return nil, &status.StatusUnauthorized
	}
	orderModel := models.OrderModel{
		ID:          uuid.NewString(),
		PaymentType: order.PaymentType,
		UserID:      userID,
	}

	userBasket, err := srv.store.Basket().GetAll(ctx, userID)
	if err != nil {
		srv.log.Error("could not load user basket",
			logs.Error(err),
			logs.String("uid", userID),
		)
		return nil, &status.StatusInternal
	}

	if len(userBasket) == 0 {
		return nil, &status.StatusBasketIsEmpty
	}

	if err := srv.store.Order().Create(ctx, orderModel); err != nil {
		srv.log.Error("could not create order", logs.Error(err))
		if errors.Is(err, storage.ErrInvalidEnumInput) {
			return nil, &status.StatusPaymentTypeInvalid
		}
		return nil, &status.StatusInternal
	}

	gotErr := false
	var errStatus *status.Status

	orderProducts := []models.OrderProductModel{}

	for _, basket := range userBasket {
		product, err := srv.store.Product().GetByID(ctx, basket.ProductID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				srv.log.Error("could not find product by id (from basket)",
					logs.String("pid", basket.ProductID),
				)
				gotErr = true
				errStatus = &status.StatusInternal
				break
			}
			srv.log.Error("could not get product by id", logs.Error(err))
			gotErr = true
			errStatus = &status.StatusInternal
			break
		}
		orderProducts = append(orderProducts, models.OrderProductModel{
			ID:           uuid.NewString(),
			OrderID:      orderModel.ID,
			ProductID:    basket.ProductID,
			Quantity:     basket.Quantity,
			ProductPrice: product.OutcomePrice,
		})
	}
	if gotErr {
		if err := srv.store.Order().Delete(ctx, orderModel.ID); err != nil {
			srv.log.Error("could not delete order", logs.Error(err))
		}
		return nil, errStatus
	}

	if err := srv.store.OrderProduct().Create(ctx, orderProducts); err != nil {
		srv.log.Error("could not create order products", logs.Error(err))

		if err := srv.store.Order().Delete(ctx, orderModel.ID); err != nil {
			srv.log.Error("could not delete order", logs.Error(err))
		}
		return nil, &status.StatusInternal
	}

	return models_v1.Response{
		Code:    200,
		Message: "Ok",
	}, nil
}

func (srv OrderService) ChangeOrderStatus(ctx context.Context, id, orderStatus string) (interface{}, *status.Status) {

	if !(orderStatus == OrderStatusFinished || orderStatus == OrderStatusCanceled) {
		return nil, &status.StatusOrderStatusInvalid
	}

	order, err := srv.store.Order().GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &status.StatusNotFound
		}
		srv.log.Error("could not get order by id", logs.Error(err),
			logs.String("oid", id))
		return nil, &status.StatusInternal
	}

	if order.DeletedAt != nil {
		return nil, &status.StatusDeleted
	}

	if order.Status != OrderStatusInProcess {
		return nil, &status.StatusNotChangable
	}

	if err := srv.store.Order().ChangeStatus(ctx, id, orderStatus); err != nil {
		if errors.Is(err, storage.ErrInvalidEnumInput) {
			return nil, &status.StatusOrderStatusInvalid
		}
		srv.log.Error("could not change order status", logs.Error(err),
			logs.String("oid", id),
			logs.String("status", orderStatus),
		)
		return nil, &status.StatusInternal
	}
	return models_v1.Response{
		Code:    200,
		Message: "Ok",
	}, nil
}

func (srv OrderService) GetByID(ctx context.Context, id string) (interface{}, *status.Status) {
	model, err := srv.store.Order().GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &status.StatusNotFound
		}
		srv.log.Error("could not get order by id", logs.Error(err),
			logs.String("oid", id))
		return nil, &status.StatusInternal
	}
	return model, nil
}

func (srv OrderService) GetAll(ctx context.Context) (interface{}, *status.Status) {
	model, err := srv.store.Order().GetAll(ctx)
	if err != nil {
		srv.log.Error("could not get order all", logs.Error(err))
		return nil, &status.StatusInternal
	}

	return model, nil
}

func (srv OrderService) GetProductByID(ctx context.Context, productID string) (interface{}, *status.Status) {
	model, err := srv.store.OrderProduct().GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &status.StatusNotFound
		}
		srv.log.Error("could not get order product by id", logs.Error(err),
			logs.String("id", productID),
		)
		return nil, &status.StatusInternal
	}
	return model, nil
}

func (srv OrderService) GetProductAll(ctx context.Context) (interface{}, *status.Status) {
	model, err := srv.store.OrderProduct().GetAll(ctx)
	if err != nil {
		srv.log.Error("could not get order all", logs.Error(err))
		return nil, &status.StatusInternal
	}

	return model, nil
}
