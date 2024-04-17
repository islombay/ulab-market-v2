package handlersv1

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	"app/pkg/helper"
	"app/pkg/logs"
	"app/storage/filestore"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"net/http"
)

// AddIconToList
// @id AddIconToList
// @router /api/icon [post]
// @summary add icon to list
// @description add icon to list
// @tags icon
// @security ApiKeyAuth
// @accept json
// @produce json
// @param add_icon_to_list formData models_v1.AddIconToList true "add icon to list"
// @param icon formData file true "icon file"
// @success 200 {object} models.IconModel "Success"
// @failure 400 {object} models_v1.Response "bad request"
// @failure 409 {object} models_v1.Response "already exists (name)"
// @failure 415 {object} models_v1.Response "invalid icon type"
// @failure 500 {object} models_v1.Response "internal error"
func (v1 *Handlers) AddIconToList(c *gin.Context) {
	var m models_v1.AddIconToList
	if c.Bind(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}
	var model = models.IconModel{
		ID:   uuid.NewString(),
		Name: m.Name,
	}

	if _, err := v1.storage.Icon().GetIconByName(context.Background(), m.Name); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get icon by name", logs.Error(err))
			return
		}
	} else {
		v1.error(c, status.StatusAlreadyExists)
		return
	}

	if _, err, msg := helper.IsValidIcon(m.Icon); err != nil {
		if errors.Is(err, helper.ErrInvalidIconType) {
			v1.error(c, status.StatusIconTypeUnkown)
			v1.log.Error("got invalid icon extension", logs.String("got", msg))
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not check file for validity", logs.Error(err))
		return
	}

	if url, err := v1.filestore.Create(m.Icon, filestore.FolderCategory, model.ID); err != nil {
		v1.error(c, status.StatusInternal)
		return
	} else {
		model.URL = url
	}

	if err := v1.storage.Icon().AddIcon(context.Background(), model); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not create icon in db list", logs.Error(err),
			logs.String("name", m.Name), logs.String("url", model.URL),
		)
		return
	}

	model.URL = v1.filestore.GetURL(model.URL)
	v1.response(c, http.StatusOK, model)
}

// GetIconsAll
// @id GetIconsAll
// @router /api/icon [get]
// @summary get all
// @description get all
// @tags icon
// @accept json
// @produce json
// @success 200 {object} []models.IconModel "Success"
// @failure 500 {object} models_v1.Response "internal error"
func (v1 *Handlers) GetIconsAll(c *gin.Context) {
	all, err := v1.storage.Icon().GetAll(context.Background())
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get all icons", logs.Error(err))
		return
	}

	for i, _ := range all {
		all[i].URL = v1.filestore.GetURL(all[i].URL)
	}

	if all == nil {
		all = []models.IconModel{}
	}

	v1.response(c, http.StatusOK, all)
}

// GetIconByID
// @id GetIconByID
// @router /api/icon/{id} [get]
// @summary get by id
// @description get by id
// @tags icon
// @accept json
// @produce json
// @param id path string true "id"
// @success 200 {object} models.IconModel "Success"
// @failure 400 {object} models_v1.Response "bad id"
// @failure 404 {object} models_v1.Response "not found"
// @failure 500 {object} models_v1.Response "internal error"
func (v1 *Handlers) GetIconByID(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	m, err := v1.storage.Icon().GetIconByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get icon by id", logs.Error(err),
			logs.String("id", id),
		)
		return
	}
	m.URL = v1.filestore.GetURL(m.URL)

	v1.response(c, http.StatusOK, m)
}

// DeleteIcon
// @id DeleteIcon
// @router /api/icon/{id} [delete]
// @summary delete icon
// @description delete icon
// @tags icon
// @accept json
// @security ApiKeyAuth
// @produce json
// @param id path string true "id"
// @success 200 {object} models_v1.Response "Success"
// @failure 400 {object} models_v1.Response "bad id"
// @failure 500 {object} models_v1.Response "internal error"
func (v1 *Handlers) DeleteIcon(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	if err := v1.storage.Icon().Delete(context.Background(), id); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not delete icon", logs.Error(err))
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}
