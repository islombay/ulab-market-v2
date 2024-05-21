package handlersv1

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/helper"
	"app/pkg/logs"
	"app/storage"
	"app/storage/filestore"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// CreateBrand godoc
// @ID createBrand
// @Router /api/brand [post]
// @Summary Create brand
// @Description Create brand
// @Tags brand
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param create_brand formData models_v1.CreateBrand true "Create brand request"
// @param image	formData file true "Brand image"
// @Success 200 {object} models.Brand "success"
// @Failure 400 {object} models_v1.Response "Bad request"
// @Failure 409 {object} models_v1.Response "Already exists"
// @Failure 500 {object} models_v1.Response "internal error"
func (v1 *Handlers) CreateBrand(c *gin.Context) {
	var m models_v1.CreateBrand
	if err := c.Bind(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		v1.log.Error("bad request", logs.Error(err))
		return
	}
	b := models.Brand{
		ID:   uuid.NewString(),
		Name: m.Name,
	}

	fmt.Println(m)
	fmt.Println(m.Image == nil)

	if m.Image != nil {
		if m.Image.Size > v1.cfg.Media.CategoryPhotoMaxSize {
			v1.error(c, status.StatusImageMaxSizeExceed)
			return
		}

		if valid, err, _ := helper.IsValidImage(m.Image); !valid && err != nil {
			if errors.Is(err, helper.ErrInvalidImageType) {
				v1.error(c, status.StatusImageTypeUnkown)
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not check the image type", logs.Error(err))
			return
		}

		url, err := v1.filestore.Create(m.Image, filestore.FolderCategory, b.ID)
		if err != nil {
			v1.log.Error("could not create brand image file", logs.Error(err))
			v1.error(c, status.StatusInternal)
			return
		}
		b.Image = &url
	}

	if err := v1.storage.Brand().Create(context.Background(), b); err != nil {
		if errors.Is(err, storage.ErrNotAffected) {
			v1.error(c, status.StatusInternal)
			v1.log.Error("got rowsAffected != 1 for creating brand", logs.Error(err))
			return
		} else if errors.Is(err, storage.ErrAlreadyExists) {
			v1.error(c, status.StatusAlreadyExists)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not create brand", logs.Error(err))
		return
	}
	v1.response(c, http.StatusOK, b)
}

// GetBrandByID
// @id getBrandById
// @router /api/brand/{id} [get]
// @tags brand
// @accept json
// @produce json
// @summary get brand by id
// @description get brand by id
// @param id path string true "brand id"
// @success 200 {object} models.Brand "brand returned"
// @failure 400 {object} models_v1.Response "Bad UUID"
// @failure 404 {object} models_v1.Response "Brand not found"
// @failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) GetBrandByID(c *gin.Context) {
	uid := c.Param("id")
	if !helper.IsValidUUID(uid) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	cat, err := v1.storage.Brand().GetByID(context.Background(), uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusBrandNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get brand by id", logs.Error(err))
		return
	}

	if cat.Image != nil {
		cat.Image = models.GetStringAddress(v1.filestore.GetURL(*cat.Image))
	}

	v1.response(c, http.StatusOK, cat)
}

// ChangeBrand
// @id changeBrand
// @router /api/brand [put]
// @summary change brand
// @description change brand name
// @tags brand
// @security ApiKeyAuth
// @accept json
// @produce json
// @param changeBrand formData models_v1.ChangeBrand true "change brand"
// @param image formData file false "image"
// @Success 200 {object} models_v1.Response "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid / no update provided"
// @Failure 404 {object} models_v1.Response "Brand not found"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) ChangeBrand(c *gin.Context) {
	var m models_v1.ChangeBrand
	if c.Bind(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	if !helper.IsValidUUID(m.ID) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	b, err := v1.storage.Brand().GetByID(context.Background(), m.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusBrandNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get brand by id", logs.Error(err), logs.String("bid", m.ID))
		return
	}

	b = &models.Brand{
		ID:    m.ID,
		Image: nil,
		Name:  "",
	}

	if m.Name == nil && m.Image == nil {
		v1.error(c, status.StatusNoUpdateProvided)
		return
	}

	if m.Name != nil {
		b.Name = *m.Name
	}

	// fmt.Println(m)

	if m.Image != nil {
		if m.Image.Size > v1.cfg.Media.CategoryPhotoMaxSize {
			v1.error(c, status.StatusImageMaxSizeExceed)
			return
		}

		if valid, err, _ := helper.IsValidImage(m.Image); !valid && err != nil {
			if errors.Is(err, helper.ErrInvalidImageType) {
				v1.error(c, status.StatusImageTypeUnkown)
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not check the image type", logs.Error(err))
			return
		}

		url, err := v1.filestore.Create(m.Image, filestore.FolderCategory, b.ID)
		if err != nil {
			v1.log.Error("could not create brand image file", logs.Error(err))
			v1.error(c, status.StatusInternal)
			return
		}
		b.Image = &url
	}

	if err := v1.storage.Brand().Change(context.Background(), *b); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not update brand", logs.Error(err), logs.String("bid", m.ID))
		return
	}
	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// GetAllBrand
// @id getAllBrand
// @router /api/brand [get]
// @tags brand
// @accept json
// @produce json
// @summary get brand all
// @description get brand
// @success 200 {object} []models.Brand "brand returned"
// @failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) GetAllBrand(c *gin.Context) {
	res, err := v1.storage.Brand().GetAll(context.Background())
	if err != nil {
		v1.log.Error("could not get all brands", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}

	for i := range res {
		if res[i].Image != nil {
			res[i].Image = models.GetStringAddress(v1.filestore.GetURL(*res[i].Image))
		}
	}
	v1.response(c, http.StatusOK, res)
}

// DeleteBrand
// @id deleteBrand
// @router /api/brand/{id} [delete]
// @tags brand
// @accept json
// @security ApiKeyAuth
// @produce json
// @summary delete brand
// @param id path string true "brand id"
// @description delete brand
// @success 200 {object} models_v1.Response "deleted successfully"
// @failure 400 {object} models_v1.Response "bad uuid"
// @failure 404 {object} models_v1.Response "brand not found"
// @failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) DeleteBrand(c *gin.Context) {
	id := c.Param("id")
	if !helper.IsValidUUID(id) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	if _, err := v1.storage.Brand().GetByID(context.Background(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusBrandNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get brand by id", logs.Error(err), logs.String("bid", id))
		return
	}
	if err := v1.storage.Brand().Delete(context.Background(), id); err != nil {
		v1.log.Error("could not delete brand", logs.Error(err), logs.String("bid", id))
		v1.error(c, status.StatusInternal)
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}
