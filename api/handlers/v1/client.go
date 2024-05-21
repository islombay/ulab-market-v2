package handlersv1

import (
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
