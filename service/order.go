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
	"strconv"
	"time"

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
	OrderStatusInProcess  = "in_process"
	OrderStatusFinished   = "finished"
	OrderStatusCanceled   = "canceled"
	OrderStatusDelivering = "delivering"
	OrderStatusPicked     = "picked"
)

func (srv OrderService) CreateOrder(ctx context.Context, order models_v1.CreateOrder) (interface{}, *status.Status) {
	userID, ok := ctx.Value(ContextUserID).(string)
	if !ok {
		srv.log.Error("invalid user id")
		return nil, &status.StatusUnauthorized
	}

	if !helper.IsValidPhone(order.ClientPhone) {
		return nil, &status.StatusBadPhone
	}

	orderModel := models.OrderModel{
		ID:              uuid.NewString(),
		OrderID:         strconv.FormatInt(time.Now().Unix(), 10),
		PaymentType:     order.PaymentType,
		UserID:          userID,
		ClientFirstName: &order.ClientFirstName,
		ClientLastName:  &order.ClientLastName,
		ClientPhone:     &order.ClientPhone,
		ClientComment:   order.ClientComment,

		DeliveryType:     order.DeliveryType,
		DeliveryAddrLat:  order.DeliveryAddrLat,
		DeliveryAddrLong: order.DeliveryAddrLong,
		DeliveryAddrName: order.DeliverAddrName,
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

		if product.Quantity < basket.Quantity {
			gotErr = true
			errStatus = &status.StatusProductQuantityTooMany
			break
		}
		orderProducts = append(orderProducts, models.OrderProductModel{
			ID:           uuid.NewString(),
			OrderID:      &orderModel.ID,
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

	for _, product_basket := range orderProducts {
		srv.store.Basket().Delete(ctx, userID, product_basket.ProductID)
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

	if orderStatusID, exists := models.OrderStatusIndexes[model.Status]; exists {
		model.StatusID = orderStatusID
	}

	orderProducts, errStatus := srv.GetOrderProducts(ctx, model.ID)
	if errStatus != nil {
		return nil, errStatus
	}

	model.Products = orderProducts.([]models.OrderProductModel)

	return model, nil
}

func (srv OrderService) GetAll(ctx context.Context, orderStatus string, pagination models.Pagination) (interface{}, *status.Status) {
	var model []models.OrderModel

	var err error
	if orderStatus == "archive" {
		model, err = srv.store.Order().GetAll(ctx, pagination, []string{"finished", "canceled"})
	} else if orderStatus == "active" {
		// model, err = srv.store.Order().GetActive(ctx)
		model, err = srv.store.Order().GetAll(ctx, pagination, []string{"in_process", "picked", "delivering"})
	} else {
		model, err = srv.store.Order().GetAll(ctx, pagination, []string{})
	}
	if err != nil {
		srv.log.Error("could not get order all archived", logs.Error(err))
		return nil, &status.StatusInternal
	}

	for i, _ := range model {
		if orderStatusID, exists := models.OrderStatusIndexes[model[i].Status]; exists {
			model[i].StatusID = orderStatusID
		}
	}

	return model, nil
}

func (srv OrderService) GetOrderProducts(ctx context.Context, order_id string) (interface{}, *status.Status) {
	model, err := srv.store.OrderProduct().GetOrderProducts(ctx, order_id)
	if err != nil {
		srv.log.Error("could not get products of order", logs.Error(err),
			logs.String("order_id", order_id))
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

func (srv OrderService) GetNewList(ctx context.Context, forCourier bool) (interface{}, *status.Status) {
	model, err := srv.store.Order().GetNew(ctx, forCourier)
	if err != nil {
		srv.log.Error("could not get new orders list", logs.Error(err))
		return nil, &status.StatusInternal
	}

	for i, _ := range model {
		if orderStatusID, exists := models.OrderStatusIndexes[model[i].Status]; exists {
			model[i].StatusID = orderStatusID
		}
	}

	return model, nil
}

func (srv OrderService) MakePicked(ctx context.Context, order_id, userID, user_type string) (interface{}, *status.Status) {
	// check the status
	// if order.status in (delivering, picked, finished, canceled)
	//		return error (not able to change)

	// check if order exists and not deleted
	model, err := srv.store.Order().GetByID(ctx, order_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &status.StatusNotFound
		}
		srv.log.Error("could not find order by id", logs.Error(err),
			logs.String("order_id", order_id))
		return nil, &status.StatusInternal
	}

	if model.DeletedAt != nil {
		return nil, &status.StatusNotFound
	}

	if model.Status == OrderStatusDelivering ||
		model.Status == OrderStatusFinished ||
		model.Status == OrderStatusCanceled {
		return nil, &status.StatusNotChangable
	}

	if user_type == "picker" {
		if model.Status == OrderStatusPicked {
			return nil, &status.StatusNotChangable
		}
		// set picked_at, picker_user_id, and status to picked
		if err := srv.store.Order().MarkPicked(ctx, order_id, userID, time.Now()); err != nil {
			srv.log.Error("could not mark order as picked", logs.Error(err),
				logs.String("order_id", order_id))
			return nil, &status.StatusInternal
		}
	} else if user_type == "courier" {
		if err := srv.store.Order().MarkPickedByCourier(ctx, order_id, userID, time.Now()); err != nil {
			srv.log.Error("could not mark order as picked by courier", logs.Error(err),
				logs.String("order_id", order_id))
			return nil, &status.StatusInternal
		}
	}

	return models_v1.Response{
		Message: "OK",
		Code:    http.StatusOK,
	}, nil
}
