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
)
