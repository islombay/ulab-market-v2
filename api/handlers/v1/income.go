package handlersv1

import (
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/helper"
	"app/pkg/logs"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateIncome	godoc
// @id			create income
// @router		/api/income [post]
// @summary		Create income
// @description Create income. Products must be given inside json body
// @tags 		income
// @security 	ApiKeyAuth
// @param		createIncome	body	models_v1.CreateIncome	true	"Create income body model"
// @success		200		{object}	models_v1.Income	"Success body"
// @failure		400		{object} 	models_v1.Response	"Bad request"
// @failure		500		{object}	models_v1.Response 	"Internal error"
func (v1 *Handlers) CreateIncome(c *gin.Context) {
	var m models_v1.CreateIncome
	if err := c.BindJSON(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, errStatus := v1.service.Income().Create(ctx, m)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, 200, resp)
}

// GetIncomeList	godoc
// @id			GetIncomeList
// @router		/api/income [get]
// @summary		Get all incomes
// @description Get all incomes list
// @tags 		income
// @security 	ApiKeyAuth
// @param		page	query		string					false	"Page number"
// @param		limit	query		string					false	"Limit of income output"
// @param		search	query		string					false	"Search inside the income"
// @success		200		{object}	models_v1.IncomeResponse	"Success body"
// @failure		400		{object} 	models_v1.Response			"Bad request"
// @failure		500		{object}	models_v1.Response 			"Internal error"
func (v1 *Handlers) GetIncomeList(c *gin.Context) {
	var m models_v1.IncomeRequest
	if err := c.Bind(&m); err != nil {
		v1.log.Error("bad request", logs.Error(err))
		v1.error(c, status.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, errStatus := v1.service.Income().GetList(ctx, m)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, resp)
}

// GetIncomeByID	godoc
// @id			GetIncomeByID
// @router		/api/income/{id} [get]
// @summary		Get income
// @description Get income
// @tags 		income
// @security 	ApiKeyAuth
// @param 		id 		path 		string 		true 	"product id"
// @success		200		{object}	models_v1.Income	"Success body"
// @failure		400		{object} 	models_v1.Response	"Bad request"
// @failure		500		{object}	models_v1.Response 	"Internal error"
func (v1 *Handlers) GetIncomeByID(c *gin.Context) {
	id := c.Param("id")
	if !helper.IsValidUUID(id) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, errStatus := v1.service.Income().GetByID(ctx, id)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, resp)
}
