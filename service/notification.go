package service

import (
	"app/pkg/logs"
	"app/service/notification"
)

type notificationService struct {
	Courier *notification.CourierNotifyService
	// Picker  *pickerNotifyService
}

func NewNotificationService(log logs.LoggerInterface) *notificationService {
	return &notificationService{
		Courier: notification.NewCourierService(log),
	}
}
