package handlersv1

import (
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/helper"
	"app/pkg/logs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateCourier		godoc
// @ID 					createCourier
// @Router 				/api/courier [POST]
// @Tags 				courier
// @Summary 			Create courier
// @Description 		Create courier
// @Accept 				json
// @Produce 			json
// @Security 			ApiKeyAuth
// @Param 		createCourier body models_v1.RegisterRequest true "Create courier body"
// @Success 	200 {object} models_v1.UUIDResponse "Successfully created"
// @Response 	400 {object} models_v1.Response "Bad request/Invalid email/Invalid phone/Invalid password"
// @Response 	401 {object} models_v1.Response "Unauthorized"
// @Response 	403 {object} models_v1.Response "Forbidden. Current user has no enough permissions to create courier"
// @Response 	409 {object} models_v1.Response "Already exists"
// @Failure 	500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) CreateCourier(c *gin.Context) {
	var m models_v1.RegisterRequest
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	resp, errStatus := v1.service.Courier().CreateCourier(m)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, resp)
}

// DeleteCourier		godoc
// @ID 					DeleteCourier
// @Router 				/api/courier/{id} [delete]
// @Tags 				courier
// @Summary 			delete courier
// @Description 		delete courier
// @Accept 				json
// @Produce 			json
// @Security 			ApiKeyAuth
// @Param 		id path string true "Courier ID"
// @Success 	200 {object} models_v1.Response
// @Failure 	400 {object} models_v1.Response "Invalid UUID"
// @Failure 	404 {object} models_v1.Response "User not found"
// @Failure 	500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) DeleteCourier(c *gin.Context) {
	uid := c.Param("id")
	if !helper.IsValidUUID(uid) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	resp, errStatus := v1.service.Courier().DeleteCourier(uid)
	if errStatus != nil {
		v1.error(c, *errStatus)
		return
	}
	v1.response(c, http.StatusOK, resp)
}

func (v1 *Handlers) CourierOrdersRealTimeConnection(c *gin.Context) {
	ws, err := v1.service.Notify().Courier.GetUpgrader(c.Writer, c.Request)
	if err != nil {
		v1.error(c, status.StatusInternal)
		return
	}

	defer ws.Close()

	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			v1.log.Debug("lost connection with client", logs.Error(err))
			return
		}
	}
}
