package handlersv1

import (
	"app/api/models"
	"app/api/status"
	auth_lib "app/pkg/auth"
	"app/pkg/logs"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
	UserIDContext       = "uid"
	TokenStaff          = "staff"
	TokenClient         = "client"
	TokenSuper          = "super"
)

func (v1 *Handlers) MiddlewareIsStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, st := v1.middlewareToken(c)
		if st != nil {
			v1.error(c, *st)
			return
		}
		if token.Type != TokenStaff {
			v1.error(c, status.StatusForbidden)
			v1.log.Debug("forbidden operation", logs.String("need", TokenSuper), logs.String("have", token.Type))
			return
		}
		c.Set(UserIDContext, token.UID)
		c.Next()
	}
}

func (v1 *Handlers) MiddlewareIsSuper() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, st := v1.middlewareToken(c)
		if st != nil {
			v1.error(c, *st)
			return
		}
		if token.Type != TokenSuper {
			v1.error(c, status.StatusForbidden)
			v1.log.Debug("forbidden operation", logs.String("need", TokenSuper), logs.String("have", token.Type))
			return
		}
		c.Set(UserIDContext, token.UID)
		c.Next()
	}
}

func (v1 *Handlers) middlewareToken(c *gin.Context) (*auth_lib.Token, *status.Status) {
	authHeader := c.Request.Header.Get(AuthorizationHeader)
	if authHeader == "" {
		return nil, &status.StatusUnauthorized
	}
	authHeaderSplit := strings.Split(authHeader, " ")
	var authToken string
	if len(authHeaderSplit) != 2 {
		//v1.error(c, status.StatusUnauthorized)
		v1.log.Debug("auth header len != 2")
		authToken = authHeader
		//return
	} else {
		authToken = authHeaderSplit[1]
	}
	token, err := auth_lib.ParseToken(authToken)
	if err != nil {
		if errors.Is(err, auth_lib.ErrTokenInvalid) {
			v1.log.Debug("got token invalid from jwt parser")
			return nil, &status.StatusUnauthorized
		} else if errors.Is(err, auth_lib.ErrTokenExpired) {
			v1.log.Debug("got token expired from jwt parser")
			return nil, &status.StatusUnauthorized
		}
		v1.log.Error("error while parsing the token", logs.Error(err))
		return nil, &status.StatusInternal
	}
	return token, nil
}

func (v1 *Handlers) MiddlewareStaffPermissionCheck(permission models.PermissionModel) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, st := v1.middlewareToken(c)
		if st != nil {
			v1.log.Debug("got status from middlewareToken", logs.Any("status", *st))
			v1.error(c, *st)
			return
		}
		userID := token.UID

		user, err := v1.storage.User().GetStaffByID(context.Background(), userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				v1.log.Debug("could not get staff by id", logs.String("uid", userID))
				v1.error(c, status.StatusUnauthorized)
				return
			}
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get staff mem by id", logs.Error(err), logs.String("uid", userID))
			return
		}

		permissionsList, err := v1.storage.Role().GetRolePermissions(context.Background(), user.RoleID)
		if err != nil {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get permissions of role", logs.Error(err))
			return
		}

		found := false
		for _, e := range permissionsList {
			if permission.ID == e.ID {
				found = true
			}
		}
		if found {
			c.Next()
		} else {
			v1.error(c, status.StatusForbidden)
		}
	}
}
