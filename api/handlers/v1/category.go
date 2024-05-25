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
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// CreateCategory godoc
// @ID createCategory
// @Router /api/category [post]
// @Summary Create category
// @Description Create category
// @Tags category
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param create_category formData models_v1.CreateCategory true "Create category request"
// @param image formData file false "Image file"
// @Success 200 {object} models_v1.ResponseID "success"
// @Failure 400 {object} models_v1.Response "Bad request"
// @failure 404 {object} models_v1.Response "Icon ID not found / Parent category not found"
// @Failure 409 {object} models_v1.Response "Already exists"
// @Failure 500 {object} models_v1.Response "internal error"
func (v1 *Handlers) CreateCategory(c *gin.Context) {
	var m models_v1.CreateCategory
	if c.Bind(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}
	if m.ParentID != "" {
		if !helper.IsValidUUID(m.ParentID) {
			v1.error(c, status.StatusBadUUID)
			return
		}
		_, err := v1.storage.Category().GetByID(context.Background(), m.ParentID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				v1.error(c, status.StatusParentCategoryNotFound)
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get by id for category", logs.Error(err))
			return
		}
	}
	var pn = models.GetStringAddress(m.ParentID)

	if _, err := v1.storage.Category().GetByName(context.Background(), m.NameRu); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get category by name", logs.Error(err))
			return
		}
	} else {
		v1.error(c, status.StatusAlreadyExists)
		return
	}

	if _, err := v1.storage.Category().GetByName(context.Background(), m.NameUz); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get category by name", logs.Error(err))
			return
		}
	} else {
		v1.error(c, status.StatusAlreadyExists)
		return
	}

	ct := models.Category{
		ID:        uuid.New().String(),
		NameUz:    m.NameUz,
		NameRu:    m.NameRu,
		ParentID:  pn,
		CreatedAt: time.Now(),
	}

	if m.IconID != nil {
		if !helper.IsValidUUID(*m.IconID) {
			v1.error(c, status.StatusBadUUID)
			return
		}
		if _, err := v1.storage.Icon().GetIconByID(context.Background(), *m.IconID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				v1.error(c, status.StatusIconNotFound)
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get icon by id", logs.Error(err))
			return
		}

		ct.IconID = m.IconID
	}

	var url string
	var err error

	if m.Image != nil {
		if m.Image.Size > v1.cfg.Media.CategoryPhotoMaxSize {
			v1.error(c, status.StatusImageMaxSizeExceed)
			v1.log.Debug("image size exceeds limit",
				logs.Any("limit", v1.cfg.Media.CategoryPhotoMaxSize),
				logs.Any("got", m.Image.Size),
			)
			return
		}
		if valid, err, contentType := helper.IsValidImage(m.Image); !valid && err != nil {
			if errors.Is(err, helper.ErrInvalidImageType) {
				v1.error(c, status.StatusImageTypeUnkown)
				v1.log.Debug("got image type", logs.String("content-type", contentType))
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not check the image type", logs.Error(err))
			return
		}
		url, err = v1.filestore.Create(m.Image, filestore.FolderCategory, ct.ID)
		if err != nil {
			v1.log.Error("could not create image file in filestore", logs.Error(err))
			v1.error(c, status.StatusInternal)
			return
		}
		ct.Image = &url
	}

	if err := v1.storage.Category().Create(context.Background(), ct); err != nil {
		if errors.Is(err, storage.ErrNotAffected) {
			v1.error(c, status.StatusInternal)
			v1.log.Error("got rowsAffected != 1 for creating category", logs.Error(err))
			return
		} else if errors.Is(err, storage.ErrAlreadyExists) {
			v1.error(c, status.StatusAlreadyExists)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not create category", logs.Error(err))
		return
	}
	v1.response(c, http.StatusOK, models_v1.ResponseID{ID: ct.ID})
}

// ChangeCategoryImage
// @ID ChangeCategoryImage
// @Router /api/category/change_image [post]
// @Summary change category image
// @Description change category image
// @Tags category
// @Accept json
// @Security ApiKeyAuth
// @Produce json
// @Param changeCategoryImage formData models_v1.ChangeCategoryImage true "change category image"
// @param image formData file false "picture file"
// @Success 200 {object} models_v1.Response "Success"
// @Failure 400 {object} models_v1.Response "Bad Request / Bad UUID / No update"
// @Failure 404 {object} models_v1.Response "Category not found / Icon not found"
// @Failure 413 {object} models_v1.Response "Image size is big"
// @Failure 415 {object} models_v1.Response "Image type is not supported"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) ChangeCategoryImage(c *gin.Context) {
	var m models_v1.ChangeCategoryImage
	if c.Bind(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}
	if !helper.IsValidUUID(m.CategoryID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	if _, err := v1.storage.Category().GetByID(context.Background(), m.CategoryID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusCategoryNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get category by id", logs.Error(err))
		return
	}

	if m.IconID != nil {
		if !helper.IsValidUUID(*m.IconID) {
			v1.error(c, status.StatusBadUUID)
			return
		}
		if _, err := v1.storage.Icon().GetIconByID(context.Background(), *m.IconID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				v1.error(c, status.StatusIconNotFound)
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get icon by id", logs.Error(err))
			return
		}
	}

	var url string = ""
	var err error

	if m.Image != nil {
		if m.Image.Size > v1.cfg.Media.CategoryPhotoMaxSize {
			v1.error(c, status.StatusImageMaxSizeExceed)
			v1.log.Debug("image size exceeds limit",
				logs.Any("limit", v1.cfg.Media.CategoryPhotoMaxSize),
				logs.Any("got", m.Image.Size),
			)
			return
		}
		if valid, err, contentType := helper.IsValidImage(m.Image); !valid && err != nil {
			if errors.Is(err, helper.ErrInvalidImageType) {
				v1.error(c, status.StatusImageTypeUnkown)
				v1.log.Debug("got image type", logs.String("content-type", contentType))
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not check the image type", logs.Error(err))
			return
		}
		url, err = v1.filestore.Create(m.Image, filestore.FolderCategory, m.CategoryID)
		if err != nil {
			v1.log.Error("could not create image file in filestore", logs.Error(err))
			v1.error(c, status.StatusInternal)
			return
		}
	}

	if err := v1.storage.Category().ChangeImage(context.Background(),
		models.GetStringAddress(m.CategoryID), models.GetStringAddress(url),
		m.IconID); err != nil {
		if errors.Is(err, storage.ErrNotAffected) {
			v1.error(c, status.StatusInternal)
			v1.log.Error("got not affected for change category image in db")
			return
		} else if errors.Is(err, storage.ErrNoUpdate) {
			v1.error(c, status.StatusNoUpdateProvided)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not change category image", logs.Error(err))
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// ChangeCategory
// @id changeCategory
// @router /api/category [put]
// @summary change category
// @description change category name and parent
// @tags category
// @security ApiKeyAuth
// @accept json
// @produce json
// @param changeCategory body models_v1.ChangeCategory true "change category. all old values must be also given"
// @Success 200 {object} models_v1.Response "Success"
// @Failure 400 {object} models_v1.Response "Bad Request / Bad UUID / Text too long"
// @Failure 404 {object} models_v1.Response "Category not found / Parent category not found / Icon not found"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) ChangeCategory(c *gin.Context) {
	var m models_v1.ChangeCategory
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	if !helper.IsValidUUID(m.ID) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	_, err := v1.storage.Category().GetByID(context.Background(), m.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusCategoryNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get category by id", logs.Error(err))
		return
	}
	if m.ParentID != "" {
		if !helper.IsValidUUID(m.ParentID) {
			v1.error(c, status.StatusBadUUID)
			return
		}
		_, err := v1.storage.Category().GetByID(context.Background(), m.ParentID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				v1.error(c, status.StatusParentCategoryNotFound)
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get by id for category", logs.Error(err))
			return
		}
	}

	if m.NameRu != "" {
		if len(m.NameRu) > 250 {
			v1.error(c, status.StatusTextTooLong)
			return
		}
	}

	if m.NameUz != "" {
		if len(m.NameUz) > 250 {
			v1.error(c, status.StatusTextTooLong)
			return
		}
	}

	var pn = models.GetStringAddress(m.ParentID)
	ct := models.Category{
		ID:       m.ID,
		NameUz:   m.NameUz,
		NameRu:   m.NameRu,
		ParentID: pn,
	}

	if models.GetStringValue(m.IconID) != "" {
		if !helper.IsValidUUID(*m.IconID) {
			v1.error(c, status.StatusBadUUID)
			return
		}
		if _, err := v1.storage.Icon().GetIconByID(context.Background(), *m.IconID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				v1.error(c, status.StatusIconNotFound)
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get icon by id", logs.Error(err))
			return
		}

		ct.IconID = m.IconID
	}

	if err := v1.storage.Category().ChangeCategory(context.Background(), ct); err != nil {
		if errors.Is(err, storage.ErrNotAffected) {
			v1.log.Error("got not affected on changing category")
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not update category", logs.Error(err), logs.String("category_id", ct.ID))
		return
	}
	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// GetCategoryByID
// @id getCategoryById
// @router /api/category/{id} [get]
// @tags category
// @accept json
// @produce json
// @summary get category by id
// @description get category by id, returns translations, and subcategories for specified category
// @param id path string true "category id"
// @success 200 {object} models.CategorySwagger "category returned"
// @failure 400 {object} models_v1.Response "Bad UUID"
// @failure 404 {object} models_v1.Response "Category not found"
// @failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) GetCategoryByID(c *gin.Context) {
	uid := c.Param("id")
	if !helper.IsValidUUID(uid) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	cat, err := v1.storage.Category().GetByID(context.Background(), uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusCategoryNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get category by id", logs.Error(err))
		return
	}

	subs, err := v1.storage.Category().GetSubcategories(context.Background(), cat.ID)
	if err != nil {
		v1.log.Error("could not load subcategories", logs.Error(err), logs.String("cid", cat.ID))
		v1.error(c, status.StatusInternal)
		return
	}
	for i, _ := range subs {
		if subs[i].Image != nil && *subs[i].Image != "" {
			subs[i].Image = models.GetStringAddress(v1.filestore.GetURL(*subs[i].Image))
		}

		if subs[i].IconID != nil && *subs[i].IconID != "" {
			ic, err := v1.storage.Icon().GetIconByID(context.Background(), *subs[i].IconID)
			if err != nil {
				v1.log.Error("could not get icon by id", logs.Error(err))
			} else {
				subs[i].IconID = models.GetStringAddress(v1.filestore.GetURL(ic.URL))
			}
		}
	}
	cat.Sub = subs

	if cat.IconID != nil && *cat.IconID != "" {
		cat.IconIDFix = *cat.IconID
		i, err := v1.storage.Icon().GetIconByID(context.Background(), *cat.IconID)
		if err != nil {
			v1.log.Error("could not get icon by id", logs.Error(err))
		} else {
			cat.IconID = models.GetStringAddress(v1.filestore.GetURL(i.URL))
		}
	}

	if cat.Image != nil && *cat.Image != "" {
		cat.Image = models.GetStringAddress(v1.filestore.GetURL(*cat.Image))
	}

	v1.response(c, http.StatusOK, cat)
}

// GetAllCategory
// @id 			getAllCategory
// @router 		/api/category [get]
// @tags 		category
// @accept 		json
// @produce 	json
// @summary 	get category all
// @param 		only_sub 	query bool	false "Only subcategory"
// @param 		limit 	query int 	false "Limit default 10"
// @param		page	query int	false "Page, default 1"
// @param		q		query string false "Query to search"
// @description get category, returns translations, and subcategories for all category
// @success 	200 {object} []models.CategorySwagger 	"category returned"
// @failure 	500 {object} models_v1.Response 		"Internal error"
func (v1 *Handlers) GetAllCategory(c *gin.Context) {
	var params models_v1.GetAllCategory
	if err := c.ShouldBind(&params); err != nil {
		v1.log.Error("bad request", logs.Error(err))
	}

	var m models.Pagination
	c.ShouldBind(&m)
	m.Fix()

	res, count, err := v1.storage.Category().GetAll(context.Background(), m, params.OnlySub)
	if err != nil {
		v1.log.Error("could not get all categories", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}

	for _, e := range res {
		if e.Image != nil && *e.Image != "" {
			e.Image = models.GetStringAddress(v1.filestore.GetURL(*e.Image))
		}
		subs, err := v1.storage.Category().GetSubcategories(context.Background(), e.ID)
		if err != nil {
			v1.log.Error("could not load subcategories", logs.Error(err), logs.String("cid", e.ID))
		}
		for i, _ := range subs {
			if subs[i].Image != nil && *subs[i].Image != "" {
				subs[i].Image = models.GetStringAddress(v1.filestore.GetURL(*subs[i].Image))
			}

			if subs[i].IconID != nil && *subs[i].IconID != "" {
				ic, err := v1.storage.Icon().GetIconByID(context.Background(), *subs[i].IconID)
				if err != nil {
					v1.log.Error("could not get icon by id", logs.Error(err))
				} else {
					subs[i].IconID = models.GetStringAddress(v1.filestore.GetURL(ic.URL))
				}
			}
		}
		if e.IconID != nil && *e.IconID != "" {
			i, err := v1.storage.Icon().GetIconByID(context.Background(), *e.IconID)
			if err != nil {
				v1.log.Error("could not get icon by id", logs.Error(err),
					logs.String("icon_id", *e.IconID))
			} else {
				e.IconID = models.GetStringAddress(v1.filestore.GetURL(i.URL))
			}
		}

		e.Sub = subs
	}
	v1.response(c, http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		Count:      count,
		Data:       res,
	})
}

// DeleteCategory
// @id deleteCategory
// @router /api/category/{id} [delete]
// @tags category
// @accept json
// @security ApiKeyAuth
// @produce json
// @summary delete category
// @param id path string true "category id"
// @description delete category & delete category translations
// @success 200 {object} models_v1.Response "deleted successfully"
// @failure 400 {object} models_v1.Response "Bad UUID"
// @failure 404 {object} models_v1.Response "Category not found"
// @failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if !helper.IsValidUUID(id) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	if _, err := v1.storage.Category().GetByID(context.Background(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusCategoryNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get category by id", logs.Error(err), logs.String("cid", id))
		return
	}
	if err := v1.storage.Category().DeleteCategory(context.Background(), id); err != nil {
		v1.log.Error("could not delete category", logs.Error(err), logs.String("cid", id))
		v1.error(c, status.StatusInternal)
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// GetCategoryBrands
// @id 		GetCategoryBrands
// @router 	/api/category/{id}/brand [get]
// @tags 	category
// @accept 	json
// @produce json
// @param 	id 	path string true 		"Category id"
// @success 200 {object} models.Brand 		"Success"
// @failure 400 {object} models_v1.Response "Bad UUID"
// @failure 500 {object} models_v1.Response "Internal server error"
func (v1 *Handlers) GetCategoryBrands(c *gin.Context) {
	id := c.Param("id")
	if !helper.IsValidUUID(id) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := v1.storage.Category().GetBrands(ctx, id)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get category brands", logs.Error(err),
			logs.String("cid", id))
		return
	}

	v1.response(c, http.StatusOK, res)
}
