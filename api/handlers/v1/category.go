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
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"net/http"
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
// @Param create_category body models_v1.CreateCategory true "Create category request"
// @Success 200 {object} models_v1.ResponseID "success"
// @Failure 400 {object} models_v1.Response "Bad request"
// @Failure 409 {object} models_v1.Response "Already exists"
// @Failure 500 {object} models_v1.Response "internal error"
func (v1 *Handlers) CreateCategory(c *gin.Context) {
	var m models_v1.CreateCategory
	if c.BindJSON(&m) != nil {
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

	ct := models.Category{
		ID:       uuid.New().String(),
		Name:     m.Name,
		ParentID: pn,
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

// AddCategoryTranslation godoc
// @ID AddCategoryTranslation
// @Router /api/category/add_translation [post]
// @Summary Create category translation
// @Description Create category translation
// @Tags category
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param create_category_translation body models_v1.CategoryTranslation true "Create category translation request"
// @Success 200 {object} models_v1.Response "success"
// @Failure 400 {object} models_v1.Response "Bad request/ Bad id"
// @Failure 404 {object} models_v1.Response "Category not found"
// @Failure 409 {object} models_v1.Response "Already exists"
// @Failure 500 {object} models_v1.Response "internal error"
func (v1 *Handlers) AddCategoryTranslation(c *gin.Context) {
	var m models_v1.CategoryTranslation
	if c.BindJSON(&m) != nil {
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
		v1.log.Error("could not get by id for category", logs.Error(err))
		return
	}
	ctm := models.CategoryTranslation{
		CategoryID:   m.CategoryID,
		Name:         m.Name,
		LanguageCode: m.LanguageCode,
	}
	if err := v1.storage.Category().AddTranslation(context.Background(), ctm); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			v1.error(c, status.StatusAlreadyExists)
			return
		} else if errors.Is(err, storage.ErrNotAffected) {
			v1.error(c, status.StatusInternal)
			v1.log.Error("got not affected on translation adding",
				logs.String("cid", m.CategoryID),
				logs.String("name", m.Name),
				logs.String("lang", m.LanguageCode),
			)
			return
		}
	}
	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
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
// @param image formData file true "picture file"
// @Success 200 {object} models_v1.Response "Success"
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid"
// @Failure 404 {object} models_v1.Response "Category not found"
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

	if _, err := v1.storage.Category().GetByID(context.Background(), m.CategoryID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusCategoryNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get category by id", logs.Error(err))
		return
	}

	url, err := v1.filestore.Create(m.Image, filestore.FolderCategory, m.CategoryID)
	if err != nil {
		v1.log.Error("could not create image file in filestore", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}

	if err := v1.storage.Category().ChangeImage(context.Background(), m.CategoryID, url); err != nil {
		if errors.Is(err, storage.ErrNotAffected) {
			v1.error(c, status.StatusInternal)
			v1.log.Error("got not affected for change category image in db")
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
// @Failure 400 {object} models_v1.Response "Bad request / bad uuid"
// @Failure 404 {object} models_v1.Response "Category not found/ parent category not found"
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
	var pn = models.GetStringAddress(m.ParentID)
	ct := models.Category{
		ID:       m.ID,
		Name:     m.Name,
		ParentID: pn,
	}
	if err := v1.storage.Category().ChangeCategory(context.Background(), ct); err != nil {
		if errors.Is(err, storage.ErrNotAffected) {
			v1.log.Error("got not affected on changing category")
		}
		v1.error(c, status.StatusInternal)
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
	translations, err := v1.storage.Category().GetTranslations(context.Background(), cat.ID)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get translations", logs.Error(err))
		return
	}
	cat.Translations = translations

	subs, err := v1.storage.Category().GetSubcategories(context.Background(), cat.ID)
	if err != nil {
		v1.log.Error("could not load subcategories", logs.Error(err), logs.String("cid", cat.ID))
		v1.error(c, status.StatusInternal)
		return
	}
	cat.Sub = subs

	if cat.Image != nil {
		cat.Image = models.GetStringAddress(v1.filestore.GetURL(*cat.Image))
	}

	v1.response(c, http.StatusOK, cat)
}

// GetAllCategory
// @id getAllCategory
// @router /api/category [get]
// @tags category
// @accept json
// @produce json
// @summary get category all
// @description get category, returns translations, and subcategories for all category
// @success 200 {object} []models.CategorySwagger "category returned"
// @failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) GetAllCategory(c *gin.Context) {
	res, err := v1.storage.Category().GetAll(context.Background())
	if err != nil {
		v1.log.Error("could not get all categories", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}

	for _, e := range res {
		tr, err := v1.storage.Category().GetTranslations(context.Background(), e.ID)
		if err != nil {
			v1.log.Error("could not get translations", logs.Error(err), logs.String("cid", e.ID))
		}
		e.Translations = tr
		if e.Image != nil {
			e.Image = models.GetStringAddress(v1.filestore.GetURL(*e.Image))
		}

		subs, err := v1.storage.Category().GetSubcategories(context.Background(), e.ID)
		if err != nil {
			v1.log.Error("could not load subcategories", logs.Error(err), logs.String("cid", e.ID))
		}
		for _, sub := range subs {
			trSub, err := v1.storage.Category().GetTranslations(context.Background(), sub.ID)
			if err != nil {
				v1.log.Error("could not get translations for sub", logs.Error(err), logs.String("cid", sub.ID))
			}
			sub.Translations = trSub

			if sub.Image != nil {
				sub.Image = models.GetStringAddress(v1.filestore.GetURL(*sub.Image))
			}
		}
		e.Sub = subs
	}
	v1.response(c, http.StatusOK, res)
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
// @failure 400 {object} models_v1.Response "bad uuid"
// @failure 404 {object} models_v1.Response "category not found"
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
