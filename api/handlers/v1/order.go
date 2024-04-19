package handlersv1

import (
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/helper"
	"app/service"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
