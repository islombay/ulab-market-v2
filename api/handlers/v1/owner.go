package handlersv1

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	auth_lib "app/pkg/auth"
	"app/pkg/helper"
	"app/pkg/logs"
	"app/pkg/start"
	"app/storage"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// CreateOwner godoc
// @ID 					createOwner
// @Router 				/api/owner [POST]
// @Tags 				owner
// @Summary 			Create owner
// @Description 		Create owner. Will return error "Already exists" if the owner is already there.
// @Accept 				json
// @Produce 			json
// @Security 			ApiKeyAuth
// @Param 				createowner body models_v1.RegisterRequest true "Create owner body"
// @Success 			200 {object} models_v1.Token 	"Successfully created"
// @Response 			400 {object} models_v1.Response "Bad request/Invalid email/Invalid phone/Invalid password"
// @Response 			401 {object} models_v1.Response "Unauthorized"
// @Response 			403 {object} models_v1.Response "Forbidden. If currect user is not super, this endpoint will return forbidden(403)"
// @Response 			409 {object} models_v1.Response "Already exists"
// @Failure 			500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) CreateOwner(c *gin.Context) {
	var m models_v1.RegisterRequest
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	usrs, err := v1.storage.User().GetStaffByRole(context.Background(), auth_lib.RoleOwner.ID)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get staff by role", logs.Error(err), logs.String("rid", auth_lib.RoleOwner.ID))
		return
	}
	if len(usrs) == 1 {
		v1.error(c, status.StatusAlreadyExists)
		v1.log.Debug("owner already exists")
		return
	} else if len(usrs) > 1 {
		v1.error(c, status.StatusAlreadyExists)
		v1.log.Error("multiple owners found in db")
		return
	}

	if !helper.IsValidEmail(m.Email) {
		v1.error(c, status.StatusBadEmail)
		return
	}
	if !helper.IsValidPhone(m.Phone) {
		v1.error(c, status.StatusBadPhone)
		return
	}
	if !helper.IsValidPassword(m.Password) {
		v1.error(c, status.StatusBadPassword)
		return
	}

	h, err := auth_lib.GetHashPassword(m.Password)
	if err != nil {
		v1.log.Error("could not generate hash password", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}

	usr := models.Staff{
		ID:          uuid.New().String(),
		Name:        m.Name,
		PhoneNumber: models.GetStringAddress(m.Phone),
		Email:       models.GetStringAddress(m.Email),
		Password:    h,
		RoleID:      auth_lib.RoleOwner.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	if err := v1.storage.User().CreateStaff(context.Background(), usr); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			v1.error(c, status.StatusAlreadyExists)
			v1.log.Debug("owner found in db while creating")
			return
		} else {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not create owner", logs.Error(err))
			return
		}
	}
	v1.response(c, http.StatusOK, models_v1.Token{Token: usr.ID})
}

// DeleteOwner godoc
// @ID deleteOwner
// @Router /api/owner/{id} [delete]
// @Tags owner
// @Summary delete owner
// @Description delete owner
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Owner ID"
// @Success 200 {object} models_v1.Response
// @Failure 400 {object} models_v1.Response "Invalid UUID"
// @Failure 404 {object} models_v1.Response "User not found"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) DeleteOwner(c *gin.Context) {
	uid := c.Param("id")
	if !helper.IsValidUUID(uid) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	if !helper.IsValidUUID(uid) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	usr, err := v1.storage.User().GetStaffByID(context.Background(), uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusUserNotFound)
			return
		}
		v1.log.Error("could not get staff by id", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}
	if usr.RoleID != auth_lib.RoleOwner.ID {
		v1.error(c, status.StatusUserNotFound)
		return
	}
	if err := v1.storage.User().DeleteStaff(context.Background(), uid); err != nil {
		if errors.Is(err, storage.ErrNotAffected) {
			v1.log.Error("rows affected is != 1", logs.String("uid", uid))
			v1.error(c, status.StatusInternal)
		} else if errors.Is(err, pgx.ErrNoRows) {
			v1.log.Error("owner not found", logs.String("uid", uid))
			v1.error(c, status.StatusUserNotFound)
		} else {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not delete owner", logs.Error(err), logs.String("uid", uid))
		}
		return
	}
	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// ChangeOwner godoc
// @ID changeOwner
// @Router /api/owner [put]
// @Summary change owner
// @Description Change owner, available for super
// @Tags owner
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param changeInfo body models_v1.ChangeAdminRequest true "Change info"
// @Success 200 {object} models_v1.Response "Success"
// @Failure 400 {object} models_v1.Response "Bad request/ invalid email/ invalid phone_number/ No update provided by user / invalid password"
// @Failure 404 {object} models_v1.Response "User/Role not found"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) ChangeOwner(c *gin.Context) {
	var m models_v1.ChangeAdminRequest
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	if !helper.IsValidUUID(m.ID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	if m.Email != "" && !helper.IsValidEmail(m.Email) {
		v1.error(c, status.StatusBadEmail)
		return
	}
	if m.Phone != "" && !helper.IsValidPhone(m.Phone) {
		v1.error(c, status.StatusBadPhone)
		return
	}
	if m.Password != "" {
		if !helper.IsValidPassword(m.Password) {
			v1.error(c, status.StatusBadPassword)
			return
		}
		pwd, err := auth_lib.GetHashPassword(m.Password)
		if err != nil {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not generate hash password", logs.Error(err))
			return
		}
		m.Password = pwd
	}
	if m.RoleID != "" {
		if _, err := v1.storage.Role().GetRoleByID(context.Background(), m.RoleID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				v1.error(c, status.StatusRoleNotFound)
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get role", logs.Error(err))
			return
		}
	}
	_, err := v1.storage.User().GetStaffByID(context.Background(), m.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusUserNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get staff by id", logs.Error(err))
		return
	}
	usr := models.Staff{
		ID:          m.ID,
		Name:        m.Name,
		Email:       &m.Email,
		PhoneNumber: &m.Phone,
		Password:    m.Password,
		RoleID:      m.RoleID,
	}
	if err := v1.storage.User().ChangeStaff(context.Background(), usr); err != nil {
		if errors.Is(err, storage.ErrNoUpdate) {
			v1.error(c, status.StatusNoUpdateProvided)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not update staff member", logs.String("uid", m.ID), logs.Error(err))
		return
	}
	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// SuperMigrateDown
// @id SuperMigrateDown
// @router /api/super/migrate-down [get]
// @security ApiKeyAuth
// @tags super
// @produce json
// @success 200 {object} models_v1.Response "success"
// @failure 500 {object} models_v1.Response "error"
func (v1 *Handlers) SuperMigrateDown(c *gin.Context) {
	if err := start.Init(&v1.cfg.DB, v1.log, true, v1.storage.Role(), v1.storage.User()); err != nil {
		v1.error(c, status.Status{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}
