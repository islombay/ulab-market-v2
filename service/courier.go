package service

import (
	"app/api/models"
	models_v1 "app/api/models/v1"
	"app/api/status"
	auth_lib "app/pkg/auth"
	"app/pkg/helper"
	"app/pkg/logs"
	"app/storage"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type courierService struct {
	store storage.StoreInterface
	log   logs.LoggerInterface
}

func NewCourierService(
	store storage.StoreInterface,
	log logs.LoggerInterface,
) *courierService {
	return &courierService{
		store: store,
		log:   log,
	}
}

func (srv *courierService) CreateCourier(m models_v1.RegisterRequest) (interface{}, *status.Status) {
	if !helper.IsValidEmail(m.Email) {
		return nil, &status.StatusBadEmail
	}
	if !helper.IsValidPhone(m.Phone) {
		return nil, &status.StatusBadPhone
	}
	if !helper.IsValidPassword(m.Password) {
		return nil, &status.StatusBadPassword
	}

	h, err := auth_lib.GetHashPassword(m.Password)
	if err != nil {
		srv.log.Error("could not generate hash password", logs.Error(err))
		return nil, &status.StatusInternal
	}

	usr := models.Staff{
		ID:          uuid.New().String(),
		Name:        m.Name,
		PhoneNumber: models.GetStringAddress(m.Phone),
		Email:       models.GetStringAddress(m.Email),
		Password:    h,
		RoleID:      auth_lib.RoleCourier.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	if err := srv.store.User().CreateStaff(context.Background(), usr); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			srv.log.Debug("admin found in db while creating")
			return nil, &status.StatusAlreadyExists
		}
		srv.log.Error("could not create admin", logs.Error(err), logs.Any("admin", usr))
		return nil, &status.StatusInternal
	}

	return models_v1.UUIDResponse{ID: usr.ID}, nil
}

func (srv *courierService) DeleteCourier(uid string) (interface{}, *status.Status) {
	usr, err := srv.store.User().GetStaffByID(context.Background(), uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &status.StatusUserNotFound
		}
		srv.log.Error("could not get staff by id", logs.Error(err))
		return nil, &status.StatusInternal
	}
	if usr.RoleID != auth_lib.RoleCourier.ID {
		return nil, &status.StatusUserNotFound
	}
	if err := srv.store.User().DeleteStaff(context.Background(), uid); err != nil {
		if errors.Is(err, storage.ErrNotAffected) {
			srv.log.Error("rows affected is != 1", logs.String("uid", uid))
			return nil, &status.StatusInternal
		} else if errors.Is(err, pgx.ErrNoRows) {
			srv.log.Error("courier not found", logs.String("uid", uid))
			return nil, &status.StatusUserNotFound
		}
		srv.log.Error("could not delete courier", logs.Error(err), logs.String("uid", uid))
		return nil, &status.StatusInternal
	}
	return models_v1.Response{
		Code:    200,
		Message: "Ok",
	}, nil
}
