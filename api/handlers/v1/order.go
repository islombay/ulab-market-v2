package handlersv1

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/helper"
	"app/service"
	"context"
	"fmt"
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
	fmt.Println(m)

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
// @param page 	query int false "Page value. Default 1"
// @param limit query int false "Limit value. Default 10"
// @success 200 {object} []models.OrderModel "Success"
// @failure 404 {object} models_v1.Response "Order not found"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) GetOrderAll(c *gin.Context) {
	status := c.Query("status")

	var m models.Pagination
	c.ShouldBind(&m)
	m.Fix()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().GetAll(ctx, status, m)
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

// GetNewOrdersList
// @id			GetNewOrdersList
// @router		/api/order/new [get]
// @summary		Get all new orders list
// @description	Get list of all new orders (especially for pickers)
// @tags 		order
// @security	ApiKeyAuth
// @accept		json
// @produce		json
// @param 		page  query int false "Page value. Default 1"
// @param 		limit query int false "Limit value. Default 10"
// @success		200	{object}	[]models.OrderModel "Success"
// @failure		500	{object}	models_v1.Response	"Internal server error"
func (v1 *Handlers) GetNewOrdersList(c *gin.Context) {

	var m models.Pagination
	c.ShouldBind(&m)
	m.Fix()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().GetNewList(ctx, m, false)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}

// GetAvailableOrdersCourier
// @id			GetAvailableOrdersCourier
// @router		/api/order/courier [get]
// @summary		Get all new orders list for couriers
// @description	Get list of all new orders (especially for couriers)
// @tags 		order
// @security	ApiKeyAuth
// @accept		json
// @produce		json
// @success		200	{object}	[]models.OrderModel "Success"
// @failure		500	{object}	models_v1.Response	"Internal server error"
func (v1 *Handlers) GetAvailableOrdersCourier(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var m models.Pagination
	c.ShouldBind(&m)
	m.Fix()

	userID, err := v1.getUserID(c)
	if err != nil {
		v1.error(c, err.(status.Status))
		return
	}

	res, errStatus := v1.service.Order().GetNewListCourier(ctx, m, userID)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}

// OrderPicked
// @id 			orderPicked
// @router		/api/order/picked/{id} [get]
// @summary		mark the order as picked (pickers)
// @description mark the order as picked
// @tags		order
// @security	ApiKeyAuth
// @param		id	path	string	true "Order id"
// @success		200	{object}	models_v1.Response	"Success"
// @failure		400	{object}	models_v1.Response	"Bad id"
// @failure		404 {object}	models_v1.Response	"Not found"
// @failure		405 {object}	models_v1.Response	"Can not change the status"
// @failure		500	{object}	models_v1.Response	"Internal server error"
func (v1 *Handlers) OrderPicked(c *gin.Context) {
	order_id := c.Param("id")
	if !helper.IsValidUUID(order_id) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	userID, err := v1.getUserID(c)
	if err != nil {
		v1.error(c, err.(status.Status))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().MakePicked(ctx, order_id, userID, "picker")
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}

// OrderMarkPickedByCourier
// @id 			OrderMarkPickedByCourier
// @router		/api/order/picked_deliver/{id} [get]
// @summary		mark the order as picked (courier)
// @description mark the order as picked for couriers (shu zakazni men olaman degan mano)
// @tags		order
// @security	ApiKeyAuth
// @param		id	path	string	true "Order id"
// @success		200	{object}	models_v1.Response	"Success"
// @failure		400	{object}	models_v1.Response	"Bad id"
// @failure		404 {object}	models_v1.Response	"Not found"
// @failure		405 {object}	models_v1.Response	"Can not change the status"
// @failure		500	{object}	models_v1.Response	"Internal server error"
func (v1 *Handlers) OrderMarkPickedByCourier(c *gin.Context) {
	order_id := c.Param("id")
	if !helper.IsValidUUID(order_id) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	userID, err := v1.getUserID(c)
	if err != nil {
		v1.error(c, err.(status.Status))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().MakePicked(ctx, order_id, userID, "courier")
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}

// ClientOrders
// @id 			ClientOrders
// @router		/api/order/myorders [get]
// @summary		get the list of orders of client (client request)
// @description get the list of orders of client (client request)
// @tags		order
// @security	ApiKeyAuth
// @param page 	query int false "Page value. Default 1"
// @param limit query int false "Limit value. Default 10"
// @success		200	{object}	models.Response		"Success"
// @failure		500	{object}	models_v1.Response	"Internal server error"
func (v1 *Handlers) ClientOrders(c *gin.Context) {
	var m models.Pagination
	c.ShouldBind(&m)
	m.Fix()

	userID, err := v1.getUserID(c)
	if err != nil {
		v1.error(c, err.(status.Status))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().GetOrderAllByClient(ctx, userID, m)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}

// OrderDelivered
// @id 				orderDelivered
// @router			/api/order/delivered/{id} [get]
// @summary		mark the order as delivered
// @description mark the order as delivered
// @tags		order
// @security	ApiKeyAuth
// @param 		id 	path 		string 		true "ID of order"
// @success		200	{object}	models.Response		"Success"
// @failure		400 {object}	models_v1.Response 	"Bad uuid"
// @failure		500	{object}	models_v1.Response	"Internal server error"
func (v1 *Handlers) OrderDelivered(c *gin.Context) {
	order_id := c.Param("id")
	if !helper.IsValidUUID(order_id) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	userID, err := v1.getUserID(c)
	if err != nil {
		v1.error(c, err.(status.Status))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().OrderDelivered(ctx, userID, order_id)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}

// CourierOrdersGetAll
// @id 			CourierOrdersGetAll
// @router		/api/order/courier/myorders [get]
// @summary		get the list of orders of courier (courier request)
// @description get the list of orders of client (courier request)
// @tags		order
// @security	ApiKeyAuth
// @param page 	query int false "Page value. Default 1"
// @param limit query int false "Limit value. Default 10"
// @success		200	{object}	models.Response		"Success"
// @failure		500	{object}	models_v1.Response	"Internal server error"
func (v1 *Handlers) CourierOrdersGetAll(c *gin.Context) {
	var m models.Pagination
	c.ShouldBind(&m)
	m.Fix()

	userID, err := v1.getUserID(c)
	if err != nil {
		v1.error(c, err.(status.Status))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().GetCourierFinishedOrders(ctx, userID, m)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}

// CourierStartDeliver
// @id			CourierStartDeliver
// @router		/api/order/courier/start [post]
// @summary		Start the delivering
// @description	Start the delivering
// @tags		order
// @security	ApiKeyAuth
// @param 		start_order body models_v1.CourierStartDeliverRequest true "Start order body"
// @success		200 {object}	models.Response		"Success"
// @failure		400 {object}	models.Response		"Bad Request"
// @failure		404 {object}	models.Response		"Not found"
// @failure		405 {object}	models.Response		"Cannot change"
// @failure		406	{object}	models.Response		"Order not yet picked \ Order not yet delivering"
// @failure		500 {object}	models.Response		"Internal server error"
func (v1 *Handlers) CourierStartDeliver(c *gin.Context) {
	var m models_v1.CourierStartDeliverRequest
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}
	userID, err := v1.getUserID(c)
	if err != nil {
		v1.error(c, err.(status.Status))
		return
	}

	m.CourierID = userID

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Order().CourierStartDelivering(ctx, m)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}
