package handlersv1

import (
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/helper"
	"app/service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateOrder
// @id CreateOrder
// @router /api/order [post]
// @summary create order
// @description create order
// @tags order
// @security ApiKeyAuth
// @accept json
// @produce json
// @param create_order body models_v1.CreateOrder true "Create order body"
// @success 200 {object} models_v1.Response "Success"
// @failure 400 {object} models_v1.Response "Bad request / payment type invalid"
// @failure 411 {object} models_v1.Response "Basket is empty"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) CreateOrder(c *gin.Context) {
	var m models_v1.CreateOrder
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	str, st := v1.getUserID(c)
	if st != nil {
		v1.error(c, st.(status.Status))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = context.WithValue(ctx, service.ContextUserID, str)

	res, errStatus := v1.service.Order().CreateOrder(ctx, m)
	if errStatus != nil {
		v1.error(c, *errStatus)
	} else {
		v1.response(c, http.StatusOK, res)
	}
}

// OrderFinish
// @id OrderFinish
// @router /api/order/finish/{id} [post]
// @summary finish the order
// @description finish the order. available only for staff members
// @tags order
// @security ApiKeyAuth
// @accept json
// @produce json
// @param id path string true "Order id"
// @success 200 {object} models_v1.Response "Success"
// @failure 400 {object} models_v1.Response "Bad request / status invalid"
// @failure 404 {object} models_v1.Response "Order not found"
// @failure 405 {object} models_v1.Response "Status already set and can not be changed"
// @failure 423 {object} models_v1.Response "Order already deleted and can not be changed"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) OrderFinish(c *gin.Context) {
	id := c.Param("id")
	if !helper.IsValidUUID(id) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*5,
	)
	defer cancel()

	res, errStatus := v1.service.Order().ChangeOrderStatus(ctx, id, "finished")
	if errStatus != nil {
		v1.error(c, *errStatus)
	} else {
		v1.response(c, http.StatusOK, res)
	}
}

// OrderCancel
// @id OrderCancel
// @router /api/order/cancel/{id} [post]
// @summary cancel the order
// @description cancel the order. available only for staff members
// @tags order
// @security ApiKeyAuth
// @accept json
// @produce json
// @param id path string true "Order id"
// @success 200 {object} models_v1.Response "Success"
// @failure 400 {object} models_v1.Response "Bad request / status invalid"
// @failure 404 {object} models_v1.Response "Order not found"
// @failure 405 {object} models_v1.Response "Status already set and can not be changed"
// @failure 423 {object} models_v1.Response "Order already deleted and can not be changed"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) OrderCancel(c *gin.Context) {
	id := c.Param("id")
	if !helper.IsValidUUID(id) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*5,
	)
	defer cancel()

	res, errStatus := v1.service.Order().ChangeOrderStatus(ctx, id, "canceled")
	if errStatus != nil {
		v1.error(c, *errStatus)
	} else {
		v1.response(c, http.StatusOK, res)
	}
}

// GetOrderByID
// @id GetOrderByID
// @router /api/order/{id} [get]
// @summary get order by id
// @description get order by id
// @tags order
// @security ApiKeyAuth
// @accept json
// @produce json
// @param id path string true "Order id"
// @success 200 {object} models.OrderModel "Success"
// @failure 400 {object} models_v1.Response "bad id"
// @failure 404 {object} models_v1.Response "Order not found"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	if !helper.IsValidUUID(id) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().GetByID(ctx, id)
	if errStatus != nil {
		v1.error(c, *errStatus)
	} else {
		v1.response(c, http.StatusOK, res)
	}
}

// GetOrderAll
// @id GetOrderAll
// @router /api/order [get]
// @summary get order all
// @description get order all
// @tags order
// @security ApiKeyAuth
// @accept json
// @produce json
// @param status query string false "Order status (active or archive)"
// @success 200 {object} []models.OrderModel "Success"
// @failure 404 {object} models_v1.Response "Order not found"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) GetOrderAll(c *gin.Context) {
	status := c.Query("status")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().GetAll(ctx, status)
	if errStatus != nil {
		v1.error(c, *errStatus)
	} else {
		v1.response(c, http.StatusOK, res)
	}
}

// GetOrderProduct
// @id GetOrderProduct
// @router /api/order/product/{id} [get]
// @summary get order product by id
// @description get order product by id
// @tags order
// @accept json
// @produce json
// @param id path string true "Order product id"
// @success 200 {object} models.OrderProductModel "Success"
// @failure 400 {object} models_v1.Response "bad id"
// @failure 404 {object} models_v1.Response "Order product not found"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) GetOrderProduct(c *gin.Context) {
	productID := c.Param("id")

	if !helper.IsValidUUID(productID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().GetProductByID(ctx, productID)
	if errStatus != nil {
		v1.error(c, *errStatus)
	} else {
		v1.response(c, http.StatusOK, res)
	}
}

// GetOrderProductAll
// @id GetOrderProductAll
// @router /api/order/product [get]
// @summary get order product all
// @description get order product all
// @tags order
// @accept json
// @produce json
// @success 200 {object} []models.OrderProductModel "Success"
// @failure 404 {object} models_v1.Response "Order not found"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) GetOrderProductAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().GetProductAll(ctx)
	if errStatus != nil {
		v1.error(c, *errStatus)
	} else {
		v1.response(c, http.StatusOK, res)
	}
}

func (v1 *Handlers) getUserID(c *gin.Context) (string, interface{}) {
	val, ok := c.Get(UserIDContext)
	if !ok {
		v1.log.Error("user-id not found")
		return "", status.StatusUnauthorized
	}
	str, ok := val.(string)
	if !ok {
		v1.log.Error("could not convert")
		return "", status.StatusUnauthorized
	}
	if !helper.IsValidUUID(str) {
		return "", status.StatusInternal
	}
	return str, nil
}

// GetArchivedOrder
// @id GetArchivedOrder
// @router /api/order/archived [get]
// @summary get order all archived
// @description get order all archived
// @tags order
// @security ApiKeyAuth
// @accept json
// @produce json
// @success 200 {object} []models.OrderModel "Success"
// @failure 404 {object} models_v1.Response "Order not found"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) GetArchivedOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().GetAllGroup(ctx, "archived")
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}

// GetActiveOrder
// @id GetActiveOrder
// @router /api/order/active [get]
// @summary get order all active
// @description get order all active
// @tags order
// @security ApiKeyAuth
// postgres://ulab_market_admin:txhTivSBM9rGMM8Z36KzZd9R8yrAvk3F@dpg-cnpt8sa1hbls738hskmg-a.oregon-postgres.render.com/ulab_market
// @accept json
// @produce json
// @success 200 {object} []models.OrderModel "Success"
// @failure 404 {object} models_v1.Response "Order not found"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) GetActiveOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().GetAllGroup(ctx, "acttive")
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}

// GetNewOrdersList
// @id			GetNewOrdersList
// @router		/api/order/new [get]
// @summary		Get all new orders list
// @description	Get list of all new orders (especially for couriers)
// @tags 		order
// @security	ApiKeyAuth
// @accept		json
// @produce		json
// @success		200	{object}	[]models.OrderModel "Success"
// @failure		500	{object}	models_v1.Response	"Internal server error"
func (v1 *Handlers) GetNewOrdersList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().GetNewList(ctx)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}
