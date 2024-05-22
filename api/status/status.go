package status

import (
	"net/http"
)

type Status struct {
	Message     string
	Code        int
	Description string
}

var (
	StatusBadRequest = Status{
		Message:     "Bad Request",
		Description: "Error request from client",
		Code:        http.StatusBadRequest,
	}
	StatusBadEmail = Status{
		Message:     "Invalid email",
		Description: "Provided invalid email",
		Code:        http.StatusBadRequest,
	}
	StatusBadGender = Status{
		Message: "Invalid gender",
		Code:    http.StatusBadRequest,
	}
	StatusBadDate = Status{
		Message: "Invalid date",
		Code:    http.StatusBadRequest,
	}
	StatusBadPhone = Status{
		Message:     "Invalid phone number",
		Description: "Provided invalid phone number",
		Code:        http.StatusBadRequest,
	}
	StatusBadPassword = Status{
		Message:     "Invalid password",
		Description: "Provided invalid password",
		Code:        http.StatusBadRequest,
	}
	StatusInternal = Status{
		Message:     "Internal server error",
		Description: "Internal server error",
		Code:        http.StatusInternalServerError,
	}

	StatusAlreadyExists = Status{
		Message:     "Already exists",
		Description: "model with specified data already exists",
		Code:        http.StatusConflict,
	}

	StatusFailedSendCode = Status{
		Message:     "Failed to send code",
		Description: "could not send verification code",
		Code:        http.StatusInternalServerError,
	}
	StatusNotImplemented = Status{
		Message:     "Not implemented",
		Description: "some function or method unimplemented",
		Code:        http.StatusNotImplemented,
	}

	StatusVerificationTypeNotFound = Status{
		Message:     "Verification type not found",
		Description: "verification type not found",
		Code:        http.StatusExpectationFailed,
	}

	StatusUserNotFound = Status{
		Message:     "User not found",
		Description: "user not found",
		Code:        http.StatusNotFound,
	}

	StatusInvalidVerificationCode = Status{
		Message:     "Invalid verification code",
		Description: "verification code is invalid",
		Code:        http.StatusNotAcceptable,
	}
	StatusUserNotVerified = Status{
		Message:     "User not verified",
		Description: "user not verified",
		Code:        http.StatusNotAcceptable,
	}
	StatusInvalidCredentials = Status{
		Message:     "Invalid credentials",
		Description: "invalid credentials is when login or password is incorrect",
		Code:        http.StatusExpectationFailed,
	}
)

var (
	StatusUnauthorized = Status{
		Message:     "Unauthorized",
		Description: "Not authorized",
		Code:        http.StatusUnauthorized,
	}
	StatusForbidden = Status{
		Message:     "Forbidden",
		Description: "Action not allowed",
		Code:        http.StatusForbidden,
	}
	StatusBadUUID = Status{
		Message:     "Bad UUID",
		Description: "Invalid UUID provided",
		Code:        http.StatusBadRequest,
	}
	StatusNoUpdateProvided = Status{
		Message:     "No update",
		Description: "user does not provide any update",
		Code:        http.StatusBadRequest,
	}
	StatusRoleNotFound = Status{
		Message:     "Role not found",
		Code:        http.StatusNotFound,
		Description: "Role not found",
	}
	StatusPermissionNotFound = Status{
		Message:     "Permission not found",
		Code:        http.StatusNotFound,
		Description: "Permission not found",
	}
	StatusParentCategoryNotFound = Status{
		Message:     "Parent category not found",
		Code:        http.StatusNotFound,
		Description: "parent category specified not found",
	}
	StatusCategoryNotFound = Status{
		Message:     "Category not found",
		Description: "category not found",
		Code:        http.StatusNotFound,
	}
	StatusImageMaxSizeExceed = Status{
		Message:     "Image size exceeds the limit",
		Description: "image size is big",
		Code:        http.StatusRequestEntityTooLarge,
	}
	StatusImageTypeUnkown = Status{
		Message:     "Invalid image type",
		Description: "image file is unknown",
		Code:        http.StatusUnsupportedMediaType,
	}
	StatusVideoTypeUnkown = Status{
		Message:     "Invalid video type",
		Description: "video file is unknown",
		Code:        http.StatusUnsupportedMediaType,
	}

	StatusIconTypeUnkown = Status{
		Message:     "Invalid icon type",
		Description: "icon file is unknown",
		Code:        http.StatusUnsupportedMediaType,
	}
)

var (
	StatusNotFound = Status{
		Message:     "Not found",
		Description: "not found",
		Code:        http.StatusNotFound,
	}

	StatusIconNotFound = Status{
		Message:     "Icon not found",
		Description: "icon not found",
		Code:        http.StatusNotFound,
	}
	StatusBasketIsEmpty = Status{
		Message:     "Basket is empty",
		Description: "Basket is empty",
		Code:        http.StatusLengthRequired,
	}
	StatusPaymentTypeInvalid = Status{
		Message:     "Payment type invalid",
		Description: "Payment type invalid",
		Code:        http.StatusBadRequest,
	}
	StatusOrderStatusInvalid = Status{
		Message:     "Order status invalid",
		Description: "Order status invalid",
		Code:        http.StatusBadRequest,
	}
	StatusDeleted = Status{
		Message:     "Deleted",
		Code:        http.StatusLocked,
		Description: "deleted",
	}
	StatusNotChangable = Status{
		Message:     "Cannot change",
		Code:        http.StatusMethodNotAllowed,
		Description: "cannot changes",
	}

	OrderNotYetPicked = Status{
		Message:     "Order not yet picked",
		Code:        http.StatusNotAcceptable,
		Description: "Order must be picked by picker, and after that courier can mark it as delivering or finished",
	}
)

var (
	StatusNameInvalid = Status{
		Message: "Name invalid",
		Code:    http.StatusBadRequest,
	}
	StatusSurnameInvalid = Status{
		Message: "Surname invalid",
		Code:    http.StatusBadRequest,
	}
	StatusTextTooLong = Status{
		Message: "Text too long",
		Code:    http.StatusBadRequest,
	}
)
