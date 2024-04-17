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
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"net/http"
)

// GetAllRoles
// @id getAllRoles
// @router /api/roles/role [get]
// @summary get all roles
// @description get all roles in db
// @security ApiKeyAuth
// @tags role
// @accept json
// @produce json
// @success 200 {object} []models.RoleModel "Roles list"
// @failure 500 {object} models_v1.Response "Internal"
func (v1 *Handlers) GetAllRoles(c *gin.Context) {
	roles, err := v1.storage.Role().GetRoles(context.Background())
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not load all roles", logs.Error(err))
		return
	}
	for _, role := range roles {
		permissionsList, err := v1.storage.Role().GetRolePermissions(context.Background(), role.ID)
		if err != nil {
			v1.error(c, status.StatusInternal)
			v1.log.Error(
				"could not get role permissions",
				logs.String("role id", role.ID),
				logs.Error(err),
			)
			return
		}
		role.Permissions = permissionsList
	}

	v1.response(c, http.StatusOK, roles)
}

// GetAllPermissions
// @id getAllPermissions
// @router /api/roles/permission [get]
// @summary get all permissions
// @description get all permissions in db
// @security ApiKeyAuth
// @tags role
// @accept json
// @produce json
// @success 200 {object} []models.PermissionModel "Permissions list"
// @failure 500 {object} models_v1.Response "Internal"
func (v1 *Handlers) GetAllPermissions(c *gin.Context) {
	permissions, err := v1.storage.Role().GetPermissions(context.Background())
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not load all permissions", logs.Error(err))
		return
	}

	v1.response(c, http.StatusOK, permissions)
}

// AttachPermissionToRole
// @id AttachPermissionToRole
// @router /api/roles/attach [post]
// @summary attach permission to role
// @description attach permission to role
// @security ApiKeyAuth
// @tags role
// @accept json
// @produce json
// @param attach_body body models_v1.AttachRoleToPermission true "attach body"
// @success 200 {object} models.AttachPermission "Attach Permission To Role"
// @failure 400 {object} models_v1.Response "Bad request / bad UUID"
// @failure 404 {object} models_v1.Response "Role not found/ Permission not found"
// @failure 409 {object} models_v1.Response "Already exists"
// @failure 500 {object} models_v1.Response "Internal"
func (v1 *Handlers) AttachPermissionToRole(c *gin.Context) {
	var m models_v1.AttachRoleToPermission
	if err := c.BindJSON(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	if !helper.IsValidUUID(m.PermissionID) || !helper.IsValidUUID(m.RoleID) {
		v1.error(c, status.StatusBadUUID)
		return
	}

	if _, err := v1.storage.Role().GetRoleByID(context.Background(), m.RoleID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusRoleNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not find role for attaching",
			logs.Error(err),
			logs.String("role id", m.RoleID),
		)
		return
	}
	if _, err := v1.storage.Role().GetPermissionByID(context.Background(), m.PermissionID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusPermissionNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not find permission for attaching",
			logs.Error(err),
			logs.String("permission id", m.RoleID),
		)
		return
	}

	ok, err := v1.storage.Role().IsRolePermissionAttachExists(context.Background(), m.RoleID, m.PermissionID)
	if err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not check for existance of attach",
			logs.Error(err),
			logs.String("role id", m.RoleID),
			logs.String("permission id", m.PermissionID),
		)
		return
	}
	if ok {
		v1.error(c, status.StatusAlreadyExists)
		return
	}

	if err := v1.storage.Role().Attach(context.Background(), m.RoleID, m.PermissionID); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			v1.error(c, status.StatusAlreadyExists)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not attach permission to role",
			logs.Error(err),
			logs.String("role id", m.RoleID),
			logs.String("permission id", m.PermissionID),
		)
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// DisAttachPermissionToRole
// @id DisAttachPermissionToRole
// @router /api/roles/attach [delete]
// @summary DisAttach Permission To Role
// @description DisAttach Permission To Role
// @security ApiKeyAuth
// @tags role
// @accept json
// @produce json
// @param disattach_body body models_v1.AttachRoleToPermission true "disattach body"
// @success 200 {object} models.AttachPermission "Disattach Permission To Role"
// @failure 400 {object} models_v1.Response "Bad request / bad UUID"
// @failure 500 {object} models_v1.Response "Internal"
func (v1 *Handlers) DisAttachPermissionToRole(c *gin.Context) {
	var m models_v1.AttachRoleToPermission
	if err := c.BindJSON(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	if !helper.IsValidUUID(m.PermissionID) || !helper.IsValidUUID(m.RoleID) {
		v1.error(c, status.StatusBadUUID)
		return
	}
	if err := v1.storage.Role().Disattach(context.Background(), m.RoleID, m.PermissionID); err != nil {
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not disattach",
			logs.Error(err),
			logs.String("role_id", m.RoleID),
			logs.String("permission_id", m.PermissionID),
		)
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    200,
		Message: "Ok",
	})
}

// CreateNewRole
// @id CreateNewRole
// @router /api/roles/role [post]
// @summary create new role
// @description create new role
// @security ApiKeyAuth
// @tags role
// @accept json
// @produce json
// @param create_new_role body models_v1.CreateNewRole true "create new role body"
// @success 200 {object} models.RoleModel "New role body"
// @failure 400 {object} models_v1.Response "Bad request"
// @failure 500 {object} models_v1.Response "Internal"
func (v1 *Handlers) CreateNewRole(c *gin.Context) {
	var m models_v1.CreateNewRole
	if err := c.BindJSON(&m); err != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	tmp := models.RoleModel{
		ID:          uuid.NewString(),
		Name:        m.Name,
		Description: models.GetStringAddress(m.Description),
		Permissions: nil,
	}

	if err := v1.storage.Role().CreateRole(context.Background(), tmp); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			v1.error(c, status.StatusAlreadyExists)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not create role", logs.Error(err))
		return
	}

	v1.response(c, http.StatusOK, tmp)
}
