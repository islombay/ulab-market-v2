package handlersv1

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	auth_lib "app/pkg/auth"
	"app/pkg/helper"
	"app/pkg/logs"
	"app/storage"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateAdmin godoc
// @ID createAdmin
// @Router /api/admin [POST]
// @Tags admin
// @Summary Create admin
// @Description Create admin.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param createAdmin body models_v1.RegisterRequest true "Create admin body"
// @Success 200 {object} models_v1.UUIDResponse "Successfully created"
// @Response 400 {object} models_v1.Response "Bad request/Invalid email/Invalid phone/Invalid password"
// @Response 401 {object} models_v1.Response "Unauthorized"
// @Response 403 {object} models_v1.Response "Forbidden. Current user has no enough permissions to create admin"
// @Response 409 {object} models_v1.Response "Already exists"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) CreateAdmin(c *gin.Context) {
	var m models_v1.RegisterRequest
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
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
		PhoneNumber: sql.NullString{Valid: true, String: m.Phone},
		Email:       sql.NullString{Valid: true, String: m.Email},
		Password:    h,
		RoleID:      auth_lib.RoleAdmin.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   sql.NullTime{Valid: false},
	}
	if err := v1.storage.User().CreateStaff(context.Background(), usr); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			v1.error(c, status.StatusAlreadyExists)
			v1.log.Debug("admin found in db while creating")
			return
		} else {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not create admin", logs.Error(err), logs.Any("admin", usr))
			return
		}
	}
	v1.response(c, http.StatusOK, models_v1.UUIDResponse{ID: usr.ID})
}

// DeleteAdmin godoc
// @ID deleteAdmin
// @Router /api/admin/{id} [delete]
// @Tags admin
// @Summary delete admin
// @Description delete admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Admin ID"
// @Success 200 {object} models_v1.Response
// @Failure 400 {object} models_v1.Response "Invalid UUID"
// @Failure 404 {object} models_v1.Response "User not found"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) DeleteAdmin(c *gin.Context) {
	uid := c.Param("id")
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
	if usr.RoleID != auth_lib.RoleAdmin.ID {
		v1.error(c, status.StatusUserNotFound)
		return
	}
	if err := v1.storage.User().DeleteStaff(context.Background(), uid); err != nil {
		if errors.Is(err, storage.ErrNotAffected) {
			v1.log.Error("rows affected is != 1", logs.String("uid", uid))
			v1.error(c, status.StatusInternal)
		} else if errors.Is(err, pgx.ErrNoRows) {
			v1.log.Error("admin not found", logs.String("uid", uid))
			v1.error(c, status.StatusUserNotFound)
		} else {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not delete admin", logs.Error(err), logs.String("uid", uid))
		}
		return
	}
	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// ChangeAdmin godoc
// @ID changeAdmin
// @Router /api/admin [put]
// @Summary change admin
// @Description Change Admin, available for owners
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param changeInfo body models_v1.ChangeAdminRequest true "Change info"
// @Success 200 {object} models_v1.Response "Success"
// @Failure 400 {object} models_v1.Response "Bad request/ invalid email/ invalid phone_number/ No update provided by user / invalid password"
// @Failure 404 {object} models_v1.Response "User/Role not found"
// @Failure 500 {object} models_v1.Response "Internal error"
func (v1 *Handlers) ChangeAdmin(c *gin.Context) {
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
		Email:       sql.NullString{Valid: m.Email != "", String: m.Email},
		PhoneNumber: sql.NullString{Valid: m.Phone != "", String: m.Phone},
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
