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

// AddToFavourite
// @id 			AddToFavourite
// @router		/api/favourite [post]
// @summary		Add product to favourite
// @description	Add product to favourite
// @tags		favourite
// @accept		json
// @security	ApiKeyAuth
// @produce		json
// @param		body body models_v1.AddToFavourite true "body"
// @success		200 {object} models.FavouriteModel 	"Success"
// @failure		400 {object} models_v1.Response		"Bad request"
// @failure		401 {object} models_v1.Response		"Unauthorized"
// @failure		404 {object} models_v1.Response		"Product not found"
// @failure		500 {object} models_v1.Response		"Internal server error"
func (v1 *Handlers) AddToFavourite(c *gin.Context) {
	var m models_v1.AddToFavourite
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	if !helper.IsValidUUID(m.ProductID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	str, st := v1.getUserID(c)
	if st != nil {
		v1.error(c, st.(status.Status))
		return
	}

	ctx, cancel := context.WithTimeout(
		context.WithValue(context.Background(), service.ContextUserID, str),
		time.Second*5,
	)
	defer cancel()

	res, errStatus := v1.service.Favourite().AddFavourite(ctx, m)
	if errStatus != nil {
		v1.error(c, *errStatus)
	} else {
		v1.response(c, http.StatusOK, res)
	}
}

func (v1 *Handlers) DeleteFromFavourite(c *gin.Context) {

}

// GetAllFavourite
// @id 			GetAllFavourite
// @router		/api/favourite [get]
// @summary		Get all from favourite
// @description	get all from favourite
// @tags		favourite
// @security	ApiKeyAuth
// @success		200 {object} []models.FavouriteModel 	"Success"
// @failure		400 {object} models_v1.Response		"Bad request"
// @failure		401 {object} models_v1.Response		"Unauthorized"
// @failure		404 {object} models_v1.Response		"Product not found"
// @failure		500 {object} models_v1.Response		"Internal server error"
func (v1 *Handlers) GetAllFavourite(c *gin.Context) {
	str, st := v1.getUserID(c)
	if st != nil {
		v1.error(c, st.(status.Status))
		return
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*5,
	)
	defer cancel()

	res, errStatus := v1.service.Favourite().GetAll(ctx, str)
	if errStatus != nil {
		v1.error(c, *errStatus)
	} else {
		v1.response(c, http.StatusOK, res)
	}
}
