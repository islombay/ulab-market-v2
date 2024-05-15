package handlersv1

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	auth_lib "app/pkg/auth"
	"app/pkg/helper"
	"app/pkg/logs"
	"app/storage"
	redis_service "app/storage/redis"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidCode = fmt.Errorf("invalid_code")
)

// ChangePassword godoc
// @ID changePassword
// @deprecatedrouter /api/auth/change_password [POST]
// @Summary ChangePassword
// @Description Change password works only for clients.
// @Tags auth
// @Accept json
// @Produce json
// @Param change_password body models_v1.ChangePassword true "Change password"
// @Success 200 {object} models_v1.Response{} "Successfully changed password"
// @Response 400 {object} models_v1.Response "Bad Request, invalid password"
// @Response 404 {object} models_v1.Response "User not found"
// @Response 406 {object} models_v1.Response "Invalid verification code"
// @Response 500 {object} models_v1.Response "Internal error"
// @Failure 501 {object} models_v1.Response "Not implemented"s
func (v1 *Handlers) ChangePassword(c *gin.Context) {
	var m models_v1.ChangePassword
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	var identityFunc func(context.Context, string) (*models.Client, error)
	if m.Type == auth_lib.VerificationEmail {
		identityFunc = v1.storage.User().GetClientByEmail
	} else {
		v1.error(c, status.StatusNotImplemented)
		return
	}
	model, err := identityFunc(context.Background(), m.Source)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusUserNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get user", logs.Error(err),
			logs.String("type", m.Type), logs.String("source", m.Source),
		)
		return
	}
	if err := v1.verifyCode(model, m.Code, true, m.Type); err != nil {
		if errors.Is(err, ErrInvalidCode) {
			v1.error(c, status.StatusInvalidVerificationCode)
			return
		}
		v1.error(c, status.StatusInternal)
		return
	}
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
	if err := v1.storage.User().ChangeStaffPassword(context.Background(), model.ID, pwd); err != nil {
		if errors.Is(err, storage.ErrNotAffected) {
			v1.log.Debug("change client password returned not affected", logs.String("uid", model.ID))
		} else {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not change password for client", logs.Error(err))
			return
		}
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    http.StatusOK,
		Message: "Ok",
	})
}

// LoginAdmin godoc
// @ID loginadmin
// @Router /api/auth/login_admin [POST]
// @Summary Login for admin
// @Description Login for admin. This is need to enter for admin panel
// @Tags auth
// @Accept json
// @Produce json
// @Param login_admin body models_v1.LoginAdminRequest true "Login request"
// @Success 200 {object} models_v1.Token "Success returning token"
// @Response 400 {object} models_v1.Response "Bad request"
// @Response 406 {object} models_v1.Response "User not verified"
// @Response 417 {object} models_v1.Response "Invalid credentials"
// @Failure 500 {object} models_v1.Response "Internal"
func (v1 *Handlers) LoginAdmin(c *gin.Context) {
	var m models_v1.LoginAdminRequest
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}
	user, err := v1.storage.User().GetStaffByLogin(context.Background(), m.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusInvalidCredentials)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get user", logs.Error(err), logs.String("login", m.Login))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(m.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			v1.log.Debug("password does not match", logs.Error(err))
			v1.error(c, status.StatusInvalidCredentials)
			return
		}
		v1.log.Debug("could not compare hash and password", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}

	memberType := TokenStaff

	r, err := v1.storage.Role().GetRoleByID(context.Background(), user.RoleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.log.Error("role specified for user not found", logs.String("rid", user.RoleID))
		} else {
			v1.log.Error("could not get role specified for user", logs.Error(err), logs.String("rid", user.RoleID))
		}
	}
	if r != nil {
		v1.log.Debug("set super role", logs.String("user's role_id", r.ID),
			logs.String("super role_id", auth_lib.RoleSuper.ID))
		if r.ID == auth_lib.RoleSuper.ID {
			memberType = TokenSuper
		}
	}

	dur := time.Hour * auth_lib.TokenExpireLife
	token, err := auth_lib.GenerateToken(user.ID, dur, memberType)
	if err != nil {
		v1.log.Error("could not generate token", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}
	v1.response(c, http.StatusOK, models_v1.Token{Token: token})
}

// LoginCourier				godoc
// @ID 						LoginCourier
// @Router 					/api/auth/login_courier [POST]
// @Summary 				Login for courier
// @Description 			Login for courier. This is need to enter for courier panel
// @Tags 					auth
// @Accept 					json
// @Produce 				json
// @Param 		login_courier body models_v1.LoginAdminRequest true "Login request"
// @Success 	200 {object} models_v1.Token "Success returning token"
// @Response 	400 {object} models_v1.Response "Bad request"
// @Response 	406 {object} models_v1.Response "User not verified"
// @Response 	417 {object} models_v1.Response "Invalid credentials"
// @Failure 	500 {object} models_v1.Response "Internal"
func (v1 *Handlers) LoginCourier(c *gin.Context) {
	var m models_v1.LoginAdminRequest
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}
	user, err := v1.storage.User().GetStaffByLogin(context.Background(), m.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusInvalidCredentials)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get user", logs.Error(err), logs.String("login", m.Login))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(m.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			v1.log.Debug("password does not match", logs.Error(err))
			v1.error(c, status.StatusInvalidCredentials)
			return
		}
		v1.log.Debug("could not compare hash and password", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}

	r, err := v1.storage.Role().GetRoleByID(context.Background(), user.RoleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.log.Error("role specified for user not found", logs.String("rid", user.RoleID))
		} else {
			v1.log.Error("could not get role specified for user", logs.Error(err), logs.String("rid", user.RoleID))
		}
	}

	if r.ID != auth_lib.RoleCourier.ID {
		v1.error(c, status.StatusInvalidCredentials)
		return
	}

	dur := time.Hour * auth_lib.TokenExpireLife
	token, err := auth_lib.GenerateToken(user.ID, dur, TokenCourier)
	if err != nil {
		v1.log.Error("could not generate token", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}
	v1.response(c, http.StatusOK, models_v1.Token{Token: token})

}

// RegisterClient godoc
// @ID register
// @deprecatedrouter /api/auth/register [POST]
// @Summary Create Client
// @Description Create Client
// @Tags auth
// @Accept json
// @Produce json
// @Param register body models_v1.RegisterRequest true "Register Request"
// @Success 200 {object} models_v1.RequestCode "Success Request that needs verification"
// @Response 400 {object} models_v1.Response{} "Bad Request / Invalid email / Invalid phone / Invalid password"
// @Response 409 {object} models_v1.Response{} "User already exists"
// @Failure 500 {object} models_v1.Response{} "Internal server error"
func (v1 *Handlers) RegisterClient(c *gin.Context) {
	var model models_v1.RegisterRequest
	if c.BindJSON(&model) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}
	if !helper.IsValidEmail(model.Email) {
		v1.error(c, status.StatusBadEmail)
		return
	}
	if !helper.IsValidPhone(model.Phone) {
		v1.error(c, status.StatusBadPhone)
		return
	}
	if !helper.IsValidPassword(model.Password) {
		v1.error(c, status.StatusBadPassword)
		return
	}

	m := models.Client{
		ID:          uuid.New().String(),
		Name:        model.Name,
		PhoneNumber: models.GetStringAddress(model.Phone),
		Email:       models.GetStringAddress(model.Email),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	if err := v1.storage.User().CreateClient(context.Background(), m); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			v1.log.Debug("client already found",
				logs.String("email", model.Email),
				logs.String("phone", model.Phone),
			)
			v1.error(c, status.StatusAlreadyExists)
			return
		}
		v1.log.Error("could not register client", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}

	verificationType := auth_lib.VerificationEmail
	verificationSource := m.Email

	oneTimeCodeDuration := 3 * time.Hour

	if err, _ := auth_lib.SendVerificationCode(*verificationSource, verificationType,
		auth_lib.CodeLength,
		v1.cache, v1.smtp,
		oneTimeCodeDuration,
	); err != nil {
		if errors.Is(err, helper.ErrInvalidEmail) {
			v1.error(c, status.StatusBadEmail)
		} else {
			v1.error(c, status.StatusFailedSendCode)
		}
		v1.log.Error("could not send verification code", logs.Error(err))
		return
	}

	v1.response(c, http.StatusOK, models_v1.RequestCode{
		Email:    *verificationSource,
		NeedCode: true,
	})
}

// VerifyCode godoc
// @ID verifyCode
// @Router /api/auth/verify_code [POST]
// @Summary Verify client's email/phone
// @Description Verify client's email/phone
// @Tags auth
// @Accept json
// @Produce json
// @Param verifyCode body models_v1.VerifyCodeRequest true "Verify code Request"
// @Success 200 {object} models_v1.Token "Success"
// @Response 400 {object} models_v1.Response "Bad Request"
// @Response 404 {object} models_v1.Response "User not found"
// @Response 406 {object} models_v1.Response "Invalid verification code"
// @Failure 500 {object} models_v1.Response "Internal server error"
// @Failure 501 {object} models_v1.Response "Not implemented email/phone verification"
func (v1 *Handlers) VerifyCode(c *gin.Context) {
	var m models_v1.VerifyCodeRequest
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	var identityFunc func(context.Context, string) (*models.Client, error)
	if m.Type == auth_lib.VerificationEmail {
		identityFunc = v1.storage.User().GetClientByEmail
	} else if m.Type == auth_lib.VerificationPhone {
		identityFunc = v1.storage.User().GetClientByPhone
	} else {
		v1.error(c, status.StatusNotImplemented)
		return
	}
	model, err := identityFunc(context.Background(), m.Source)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			v1.error(c, status.StatusUserNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get user", logs.Error(err),
			logs.String("type", m.Type), logs.String("source", m.Source),
		)
		return
	}
	if err := v1.verifyCode(model, m.Code, false, m.Type); err != nil {
		if errors.Is(err, ErrInvalidCode) {
			v1.error(c, status.StatusInvalidVerificationCode)
			return
		}
		v1.error(c, status.StatusInternal)
		return
	}

	dur := time.Hour * auth_lib.TokenExpireLife
	token, err := auth_lib.GenerateToken(model.ID, dur, TokenClient)
	if err != nil {
		v1.log.Error("could not generate token", logs.Error(err))
		v1.error(c, status.StatusInternal)
		return
	}
	v1.response(c, http.StatusOK, models_v1.Token{Token: token})
}

func (v1 *Handlers) verifyCode(model *models.Client, reqCode string, deleteAfterCheck bool, vType string) error {
	var src string
	if vType == auth_lib.VerificationEmail {
		src = *model.Email
	} else if vType == auth_lib.VerificationPhone {
		src = *model.PhoneNumber
	}
	code, exp, err := v1.cache.Code().GetCode(context.Background(), src)
	if err != nil {
		if errors.Is(err, redis_service.ErrKeyNotFound) {
			v1.log.Debug("verification code key not found")
			return ErrInvalidCode
		}
		v1.log.Error("could not get code in cache", logs.Error(err))
		return err
	}

	if exp.Before(time.Now()) {
		v1.log.Debug("verification code already expired")
		return ErrInvalidCode
	}

	if code != reqCode {
		v1.log.Debug("verification code does not match")
		return ErrInvalidCode
	}

	if deleteAfterCheck {
		v1.log.Debug("deleting verification code", logs.String("code", code))
		v1.cache.Code().DeleteCode(context.Background(), *model.Email)
	}
	return nil
}

// Login godoc
// @ID login
// @Router /api/auth/login [POST]
// @Summary Login
// @Description Login
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models_v1.LoginRequest true "Login request"
// @Success 200 {object} models_v1.RequestCode "Success. Now needs to verify verification code /api/verify_code"
// @Response 400 {object} models_v1.Response "Bad request / Bad Email / Bad Phone"
// @Response 417 {object} models_v1.Response "Invalid type"
// @Failure 500 {object} models_v1.Response "Internal"
// @failure 501 {object} models_v1.Response "Not implemented (phone verification)"
func (v1 *Handlers) Login(c *gin.Context) {
	var m models_v1.LoginRequest
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	var identityFunc func(ctx context.Context, l string) (*models.Client, error)
	// if m.Type == auth_lib.VerificationEmail {
	// 	if !helper.IsValidEmail(m.Source) {
	// 		v1.error(c, status.StatusBadEmail)
	// 		return
	// 	}
	// 	identityFunc = v1.storage.User().GetClientByEmail
	// } else
	if m.Type == auth_lib.VerificationPhone {
		if !helper.IsValidPhone(m.Source) {
			v1.error(c, status.StatusBadPhone)
			return
		}
		identityFunc = v1.storage.User().GetClientByPhone
	} else {
		v1.error(c, status.StatusVerificationTypeNotFound)
		return
	}

	user, err := identityFunc(context.Background(), m.Source)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			user = &models.Client{
				ID:        uuid.NewString(),
				Name:      "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: nil,
			}

			if m.Type == auth_lib.VerificationPhone {
				user.PhoneNumber = &m.Source
			} else if m.Type == auth_lib.VerificationEmail {
				user.Email = &m.Source
			}
			if err := v1.storage.User().CreateClient(context.Background(), *user); err != nil {
				v1.log.Error("could not register client", logs.Error(err))
				v1.error(c, status.StatusInternal)
				return
			}
		} else {
			v1.error(c, status.StatusInternal)
			v1.log.Error("could not get user", logs.Error(err), logs.String("login", m.Source))
			return
		}
	}

	verificationType := m.Type
	verificationSource := m.Source

	oneTimeCodeDuration := time.Hour
	oneTimeCode := ""
	if err, code := auth_lib.SendVerificationCode(verificationSource, verificationType,
		auth_lib.CodeLength,
		v1.cache, v1.smtp,
		oneTimeCodeDuration,
	); err != nil {
		if errors.Is(err, helper.ErrInvalidEmail) {
			v1.error(c, status.StatusBadEmail)
		} else {
			v1.error(c, status.StatusFailedSendCode)
		}
		v1.log.Error("could not send verification code",
			logs.Error(err),
			logs.String("source", verificationSource),
		)
		return
	} else {
		oneTimeCode = code
	}

	res := models_v1.RequestCode{
		NeedCode: true,
	}
	if verificationType == auth_lib.VerificationPhone {
		res.Phone = verificationSource
		res.Code = oneTimeCode
	} else if verificationType == auth_lib.VerificationEmail {
		res.Email = verificationSource
	}

	v1.response(c, http.StatusOK, res)
}

// RequestCode godoc
// @ID requestCode
// @deprecatedrouter /api/auth/request_code [POST]
// @Summary Request code
// @Description request code is needed when password is forgotten
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models_v1.RequestCodeRequest true "Request code request"
// @Success 200 {object} models_v1.Response "Success returning token"
// @Response 400 {object} models_v1.Response "Bad request / Invalid email / Invalid phone"
// @Response 404 {object} models_v1.Response "User not found"
// @Failure 500 {object} models_v1.Response "Internal / Failed to send code"
// @Failure 501 {object} models_v1.Response "Not implemented email/phone verification"
func (v1 *Handlers) RequestCode(c *gin.Context) {
	var m models_v1.RequestCodeRequest
	if c.BindJSON(&m) != nil {
		v1.error(c, status.StatusBadRequest)
		return
	}

	var identityFunc func(context.Context, string) (*models.Client, error)
	if m.Type == auth_lib.VerificationEmail {
		if !helper.IsValidEmail(m.Source) {
			v1.error(c, status.StatusBadEmail)
			return
		}

		identityFunc = v1.storage.User().GetClientByEmail
	} else {
		v1.error(c, status.StatusNotImplemented)
		return
	}

	_, err := identityFunc(context.Background(), m.Source)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			v1.error(c, status.StatusUserNotFound)
			return
		}
		v1.error(c, status.StatusInternal)
		v1.log.Error("could not get user", logs.Error(err),
			logs.String("type", m.Type), logs.String("source", m.Source),
		)
		return
	}

	verificationType := m.Type
	verificationSource := m.Source

	oneTimeCodeDuration := 3 * time.Hour

	if err, _ := auth_lib.SendVerificationCode(verificationSource, verificationType,
		auth_lib.CodeLength,
		v1.cache, v1.smtp,
		oneTimeCodeDuration,
	); err != nil {
		if errors.Is(err, helper.ErrInvalidEmail) {
			v1.error(c, status.StatusBadEmail)
		} else {
			v1.error(c, status.StatusFailedSendCode)
		}
		v1.log.Error("could not send verification code", logs.Error(err))
		return
	}

	v1.response(c, http.StatusOK, models_v1.Response{
		Code:    http.StatusOK,
		Message: "Ok",
	})
}
