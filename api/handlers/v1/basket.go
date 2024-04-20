package handlersv1

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/helper"
	"app/pkg/logs"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"net/http"
	"time"
)

// AddToBasket
// @id AddToBasket
// @router /api/basket [post]
// @summary add product to basket
// @description add product to basket
// @tags basket
// @security ApiKeyAuth
// @accept json
// @produce json
// @param add_to_basket body models_v1.AddToBasket true "Add product to basket"
// @success 200 {object} models_v1.Response "success"
// @failure 400 {object} models_v1.Response "bad request / bad uuid"
// @failure 409 {object} models_v1.Response "already found"
// @failure 500 {object} models_v1.Response "internal error"
func (v1 *Handlers) AddToBasket(c *gin.Context) {
	var m models_v1.AddToBasket
	if err := c.BindJSON(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	if !helper.IsValidUUID(m.ProductID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	val, ok := c.Get(UserIDContext)
	if !ok {
		v1.error(c, status.StatusUnauthorized)
		return
	}

	str, ok := val.(string)
	if !ok {
		v1.error(c, status.StatusUnauthorized)
		return
	}

	if !helper.IsValidUUID(str) {
		v1.error(c, status.StatusInternal)
		return
	}

	if _, err := v1.storage.Basket().Get(context.Background(),
		str,
		m.ProductID,
	); err == nil {
		v1.error(c, status.StatusAlreadyExists)
		return
	} else if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get basket product", logs.Error(err),
				logs.String("uid", str), logs.String("pid", m.ProductID))
			return
		}
	}

	if err := v1.storage.Basket().Add(
		context.Background(),
		str,
		m.ProductID,
		m.Quantity,
		time.Now(),
	); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not add product to basket", logs.Error(err))
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// GetBasket
// @id GetBasket
// @router /api/basket [get]
// @summary get all products from basket
// @description get all products from basket
// @tags basket
// @security ApiKeyAuth
// @accept json
// @produce json
// @success 200 {object} models_v1.GetBasket "success"
// @failure 500 {object} models_v1.Response "internal error"
func (v1 *Handlers) GetBasket(c *gin.Context) {
	val, ok := c.Get(UserIDContext)
	if !ok {
		v1.error(c, status.StatusUnauthorized)
		return
	}

	str, ok := val.(string)
	if !ok {
		v1.error(c, status.StatusUnauthorized)
		return
	}

	if !helper.IsValidUUID(str) {
		v1.error(c, status.StatusInternal)
		return
	}

	res, err := v1.storage.Basket().GetAll(context.Background(), str)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get all baskets for user", logs.Error(err),
			logs.String("uid", str),
		)
		return
	}
	if res == nil {
		res = []models.BasketModel{}
	}

	resBody := models_v1.GetBasket{}

	for _, e := range res {
		product, err := v1.storage.Product().GetByID(context.Background(), e.ProductID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				v1.error(c, status.StatusInternal)
				v1.log.Error("product from basket not found",
					logs.String("product_id", e.ProductID),
				)
				return
			}
		}
		tmp := models_v1.GetBasketProduct{
			ID:       product.ID,
			NameRu:   product.NameRu,
			NameUz:   product.NameUz,
			Price:    product.OutcomePrice,
			Quantity: e.Quantity,
		}
		if product.MainImage != nil {
			tmp.MainImage = models.GetStringAddress(v1.filestore.GetURL(*product.MainImage))
		}

		resBody.Products = append(resBody.Products, tmp)
		resBody.TotalPrice += tmp.Price * float64(tmp.Quantity)
	}
	v1.response(c, http.StatusOK, resBody)
}

// DeleteFromBasket
// @id DeleteFromBasket
// @router /api/basket [delete]
// @summary delete from basket
// @description delete from basket
// @tags basket
// @accept json
// @produce json
// @security ApiKeyAuth
// @param delete_from_basket body models_v1.RemoveFromBasket true "Remove from basket information"
// @success 200 {object} models_v1.Response "Success"
// @failure 400 {object} models_v1.Response "Bad request / invalid product id"
// @failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) DeleteFromBasket(c *gin.Context) {
	var m models_v1.RemoveFromBasket
	if err := c.BindJSON(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	if !helper.IsValidUUID(m.ProductID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	val, ok := c.Get(UserIDContext)
	if !ok {
		v1.error(c, status.StatusUnauthorized)
		return
	}

	str, ok := val.(string)
	if !ok {
		v1.error(c, status.StatusUnauthorized)
		return
	}

	if !helper.IsValidUUID(str) {
		v1.error(c, status.StatusInternal)
		return
	}

	if err := v1.storage.Basket().Delete(context.Background(), str, m.ProductID); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not delete from basket",
			logs.Error(err),
			logs.String("uid", str),
			logs.String("pid", m.ProductID),
		)
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// ChangeBasket
// @id ChangeBasket
// @router /api/basket [put]
// @tags basket
// @security ApiKeyAuth
// @accept json
// @produce json
// @param change_basket body models_v1.ChangeBasket true "Change basket body"
// @success 200 {object} models_v1.Response "Success"
// @failure 400 {object} models_v1.Response "Bad request / bad quantity / bad id (product)"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) ChangeBasket(c *gin.Context) {
	var m models_v1.ChangeBasket
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	if m.Quantity == 0 {
		v1.error(c, status.StatusBadRequest)
		return
	}

	if !helper.IsValidUUID(m.ProductID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	val, ok := c.Get(UserIDContext)
	if !ok {
		v1.log.Error("user-id not found")
		v1.error(c, status.StatusUnauthorized)
		return
	}

	str, ok := val.(string)
	if !ok {
		v1.log.Error("could not convert")
		v1.error(c, status.StatusUnauthorized)
		return
	}

	if !helper.IsValidUUID(str) {
		v1.error(c, status.StatusInternal)
		return
	}

	// TODO: check whether the quantity exists

	if err := v1.storage.Basket().ChangeQuantity(
		context.Background(),
		m.ProductID,
		str,
		m.Quantity,
	); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not update quantity of product in basket",
			logs.Error(err), logs.String("uid", str),
			logs.String("pid", m.ProductID), logs.Int("q", m.Quantity),
		)
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}
