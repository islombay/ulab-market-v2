package handlersv1

import (
	"app/api/models"
	"app/api/status"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetClientList
// @id			GetClientList
// @router		/api/client [get]
// @security	ApiKeyAuth
// @tags 		client
// @summary		Get list of clients ( only staff )
// @description	Get list of clients ( all information )
// @success 	200	{object}	[]models.ClientSwagger	"List of clients"
// @failure		500 {object}	models_v1.Response		"Internal server error"
func (v1 *Handlers) GetClientList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Client().GetList(ctx)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}

// ClientGetMe
// @id			ClientGetMe
// @router		/api/client/getme [get]
// @security	ApiKeyAuth
// @tags 		client
// @summary		get client information (self)
// @description	get client information (self)
// @success 	200	{object}	models.ClientSwagger	"List of clients"
// @failure		500 {object}	models_v1.Response		"Internal server error"
func (v1 *Handlers) ClientGetMe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	userID, err := v1.getUserID(c)
	if err != nil {
		v1.error(c, err.(status.Status))
		return
	}

	resp, errStatus := v1.service.Client().GetMe(ctx, userID)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}

	v1.response(c, http.StatusOK, resp)
}

// ClientUpdate
// @id			ClientUpdate
// @router		/api/client/me [put]
// @security 	ApiKeyAuth
// @tags 		client
// @summary 	update client (client)
// @description update client. only for clients
// @param		body body 	models.ClientUpdate	true "Update body"
// @success		200 {object}	models_v1.Response "Success"
// @failure		400 {object}	models_v1.Response "Bad email, bad gender, bad birthdate, bad name/surname"
// @failure		409	{object}	models_v1.Response "Email already exists"
// @failure		500	{object} 	models_v1.Response "internal server error"
func (v1 *Handlers) ClientUpdate(c *gin.Context) {
	var m models.ClientUpdate
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	userID, err := v1.getUserID(c)
	if err != nil {
		v1.error(c, err.(status.Status))
		return
	}

	m.ID = userID

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, errStatus := v1.service.Client().Update(ctx, m)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, res)
}
