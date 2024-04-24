package handlersv1

import (
	models_v1 "app/api/models/v1"
	"app/pkg/logs"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateStorage    godoc
// @Router       /api/storage [POST]
// @Summary      Create a new storage
// @Description  Create a new storage
// @Tags         storage
// @security ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        storage  body models_v1.CreateStorage false "storage"
// @Success      201  {object}  models_v1.Storage
// @Failure      400  {object}  models_v1.Response "Bad request"
// @Failure      404  {object}  models_v1.Response "Page not found"
// @Failure      500  {object}  models_v1.Response "Interval server error"
func (v1 *Handlers) CreateStorage(c *gin.Context) {
	fmt.Println("here")
	request := models_v1.CreateStorage{}
	if err := c.ShouldBindJSON(&request); err != nil {
		v1.log.Error("error is while reading body", logs.Error(err))
		return
	}

	fmt.Println("req", request)
	resp, err := v1.service.Store().CreateStore(context.Background(), request)
	if err != nil {
		v1.log.Error("error is while creating  storage", logs.Error(err))
		return
	}
	handleResponse(c, v1.log, "success", http.StatusOK, resp)
}

// GetStorageByID     godoc
// @Router       /api/storage/{id} [GET]
// @Summary      Get storage  by id
// @Description  get storage  by id
// @Tags         storage
// @security ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "storage_id"
// @Success      201  {object}  models_v1.Storage
// @Failure      400  {object}  models_v1.Response "Bad request"
// @Failure      404  {object}  models_v1.Response "Page not found"
// @Failure      500  {object}  models_v1.Response "Interval server error"
func (v1 *Handlers) GetStorageByID(c *gin.Context) {
	id := c.Param("id")

	resp, err := v1.service.Store().GetStoreByID(context.Background(), id)
	if err != nil {
		handleResponse(c, v1.log, "error is while getting storage by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, v1.log, "success", http.StatusOK, resp)
}

// GetStorageList godoc
// @Router       /api/storage  [GET]
// @Summary      Get storage  list
// @Description  get storage list
// @Tags         storage
// @security ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      201  {object}  models_v1.StorageResponse
// @Failure      400  {object}  models_v1.Response "Bad request"
// @Failure      404  {object}  models_v1.Response "Page not found"
// @Failure      500  {object}  models_v1.Response "Interval server error"
func (v1 *Handlers) GetStorageList(c *gin.Context) {
	var (
		page, limit int
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, v1.log, "error is while converting page", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, v1.log, "error is while converting limit", http.StatusBadRequest, err.Error())
		return
	}

	resp, err := v1.service.Store().GetStoreList(context.Background(), models_v1.StorageRequest{
		Page:   page,
		Limit:  limit,
		Search: c.Query("search"),
	})

	if err != nil {
		handleResponse(c, v1.log, "error is while get storage list", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, v1.log, "success", http.StatusOK, resp)
}

// UpdateStorage godoc
// @Router       /api/storage/{id} [PUT]
// @Summary      Update storage
// @Description  update storage
// @Tags         storage
// @security ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "storage"
// @Param        sale_point  body models_v1.UpdateStorage false "storage"
// @Success      201  {object}  models_v1.Storage
// @Failure      400  {object}  models_v1.Response "Bad request"
// @Failure      404  {object}  models_v1.Response "Page not found"
// @Failure      500  {object}  models_v1.Response "Interval server error"
func (v1 *Handlers) UpdateStorage(c *gin.Context) {
	request := models_v1.UpdateStorage{}
	uid := c.Param("id")

	if err := c.ShouldBindJSON(&request); err != nil {
		handleResponse(c, v1.log, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	request.ID = uid

	resp, err := v1.service.Store().UpdateStore(context.Background(), request)
	if err != nil {
		handleResponse(c, v1.log, "error is while getting storage by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, v1.log, "updated", http.StatusOK, resp)
}

// DeleteStorage  godoc
// @Router       /api/storage/{id} [DELETE]
// @Summary      Delete storage
// @Description  delete storage
// @Tags         storage
// @security ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "storage_id"
// @Success      201  {object}  models.Response
// @Failure      400  {object}  models_v1.Response "Bad request"
// @Failure      404  {object}  models_v1.Response "Page not found"
// @Failure      500  {object}  models_v1.Response "Interval server error"
func (v1 *Handlers) DeleteStorage(c *gin.Context) {
	uid := c.Param("id")

	if err := v1.service.Store().DeleteStore(context.Background(), uid); err != nil {
		handleResponse(c, v1.log, "error is while delete storage", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, v1.log, "deleted", http.StatusOK, "storage deleted!")
}
