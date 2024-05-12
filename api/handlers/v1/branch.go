package handlersv1

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/helper"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// AddBranch godoc
// @ID AddBranch
// @Router /api/branch [post]
// @Summary Create branch
// @Description Create branch
// @Tags branch
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param create_brand body models_v1.CreateBranch true "Create branch request"
// @Success 200 {object} models.BranchModel "success"
// @Failure 400 {object} models_v1.Response "Bad request"
// @Failure 409 {object} models_v1.Response "Already exists"
// @Failure 500 {object} models_v1.Response "internal error"
func (v1 *Handlers) AddBranch(c *gin.Context) {
	v1.error(c, status.Status{
		Message: "Not found",
		Code:    http.StatusNotFound,
	})
	return
	var m models_v1.CreateBranch
	if err := c.BindJSON(&m); err != nil {
		v1.log.Debug("got bad request for create branch", logs.Error(err))
		v1.error(c, status.StatusBadRequest)
		return
	}

	// TODO: implement open_time and close_time checking
	b := models.BranchModel{
		ID:   uuid.NewString(),
		Name: m.Name,

		OpenTime:  m.OpenTime,
		CloseTime: m.CloseTime,
	}
	if err := v1.storage.Branch().Create(context.Background(), b); err != nil {
		if errors.Is(err, storage.ErrNotAffected) {
			v1.error(c, status.StatusInternal)
			v1.log.Error("got rowsAffected != 1 for creating branch", logs.Error(err))
			return
		} else if errors.Is(err, storage.ErrAlreadyExists) {
			v1.error(c, status.StatusAlreadyExists)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not create branch", logs.Error(err))
		return
	}
	v1.response(c, http.StatusOK, b)
}

// GetBranchByID
// @id GetBranchByID
// @router /api/branch/{id} [get]
// @tags branch
// @accept json
// @produce json
// @summary get branch by id
// @description get branch by id
// @param id path string true "branch id"
// @success 200 {object} models.BranchModel "branch returned"
// @failure 400 {object} models_v1.Response "Bad UUID"
// @failure 404 {object} models_v1.Response "Branch not found"
// @failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) GetBranchByID(c *gin.Context) {
	uid := c.Param("id")
	if !helper.IsValidUUID(uid) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	cat, err := v1.storage.Branch().GetByID(context.Background(), uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get branch by id", logs.Error(err))
		return
	}

	v1.response(c, http.StatusOK, cat)
}

// ChangeBranch
// @id ChangeBranch
// @router /api/branch [put]
// @summary change branch
// @description change branch name
// @tags branch
// @security ApiKeyAuth
// @accept json
// @produce json
// @param changeBrand body models_v1.ChangeBranch true "change branch"
// @Success 200 {object} models_v1.Response "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid / no update provided"
// @Failure 404 {object} models_v1.Response "Brand not found"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) ChangeBranch(c *gin.Context) {
	var m models_v1.ChangeBranch
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	if !helper.IsValidUUID(m.ID) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	b, err := v1.storage.Branch().GetByID(context.Background(), m.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get branch by id", logs.Error(err), logs.String("bid", m.ID))
		return
	}
	if b.Name == m.Name {
		v1.error(c, status.StatusNoUpdateProvided)
		return
	}
	b = &models.BranchModel{
		ID:        m.ID,
		Name:      m.Name,
		OpenTime:  m.OpenTime,
		CloseTime: m.CloseTime,
	}
	if err := v1.storage.Branch().Change(context.Background(), *b); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not update branch", logs.Error(err), logs.String("bid", m.ID))
		return
	}
	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// GetAllBranches
// @id GetAllBranches
// @router /api/branch [get]
// @tags branch
// @accept json
// @produce json
// @summary get branch all
// @description get branch
// @success 200 {object} []models.BranchModel "branch returned"
// @failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) GetAllBranches(c *gin.Context) {
	res, err := v1.storage.Branch().GetAll(context.Background())
	if err != nil {
		v1.log.Error("could not get all branches", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}
	v1.response(c, http.StatusOK, res)
}

// DeleteBranch
// @id DeleteBranch
// @router /api/branch/{id} [delete]
// @tags branch
// @accept json
// @security ApiKeyAuth
// @produce json
// @summary delete branch
// @param id path string true "branch id"
// @description delete branch
// @success 200 {object} models_v1.Response "deleted successfully"
// @failure 400 {object} models_v1.Response "bad uuid"
// @failure 404 {object} models_v1.Response "brand not found"
// @failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) DeleteBranch(c *gin.Context) {
	id := c.Param("id")
	if !helper.IsValidUUID(id) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	if _, err := v1.storage.Branch().GetByID(context.Background(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusBrandNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get branch by id", logs.Error(err), logs.String("bid", id))
		return
	}
	if err := v1.storage.Branch().Delete(context.Background(), id); err != nil {
		v1.log.Error("could not delete branch", logs.Error(err), logs.String("bid", id))
		v1.error(c, status.StatusInternal)
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}
